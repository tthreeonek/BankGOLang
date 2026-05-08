package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"bank-api/internal/config"
	"bank-api/internal/handler"
	"bank-api/internal/repository/postgres"
	"bank-api/internal/router"
	"bank-api/internal/service"
	"bank-api/pkg/crypto"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/openpgp"
)

var pgpEntity *openpgp.Entity // глобальный ключ PGP для шифрования (в реальном проекте загружается из конфига)

func initPGP() error {
	// Генерация новой пары ключей (упрощённо, без сохранения)
	entity, err := openpgp.NewEntity("bank", "card-encryption", "bank@local", nil)
	if err != nil {
		return err
	}
	pgpEntity = entity
	return nil
}

// Функция-обёртка для PGP-шифрования
func encryptWithPGP(plaintext string) (string, error) {
	if pgpEntity == nil {
		return "", fmt.Errorf("PGP entity not initialized")
	}
	return crypto.EncryptPGP(plaintext, pgpEntity)
}

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Инициализация PGP
	if err := initPGP(); err != nil {
		log.Fatal("PGP init failed:", err)
	}

	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode))
	if err != nil {
		logrus.Fatal("DB open error:", err)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepo(db)
	accountRepo := postgres.NewAccountRepo(db)
	cardRepo := postgres.NewCardRepo(db)
	creditRepo := postgres.NewCreditRepo(db)

	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	accountSvc := service.NewAccountService(accountRepo)
	// Передаём функцию PGP-шифрования
	cardSvc := service.NewCardService(cardRepo, accountRepo, encryptWithPGP, cfg.HMACSecret)
	transferSvc := service.NewTransferService(accountRepo)
	cbrSvc := service.NewCBRService()
	creditSvc := service.NewCreditService(creditRepo, cbrSvc)
	emailSvc := service.NewEmailService(cfg.SMTPHost, 587, cfg.SMTPUser, cfg.SMTPPass)
	analyticsSvc := service.NewAnalyticsService()

	authH := handler.NewAuthHandler(authSvc)
	accountH := handler.NewAccountHandler(accountSvc)
	cardH := handler.NewCardHandler(cardSvc)
	transferH := handler.NewTransferHandler(transferSvc)
	creditH := handler.NewCreditHandler(creditSvc)
	analyticsH := handler.NewAnalyticsHandler(analyticsSvc)

	r := router.NewRouter(authH, accountH, cardH, transferH, creditH, analyticsH, cfg.JWTSecret)

	// Шедулер кредитных платежей
	go func() {
		ticker := time.NewTicker(12 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				processOverduePayments(db, creditRepo, accountRepo, userRepo, emailSvc)
			}
		}
	}()

	logrus.Info("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func processOverduePayments(db *sql.DB, cr *postgres.CreditRepo, ar *postgres.AccountRepo, ur *postgres.UserRepo, es *service.EmailService) {
	ctx := context.Background()
	now := time.Now()
	overdue, err := cr.GetOverduePayments(ctx, now)
	if err != nil {
		logrus.Error("Overdue fetch error:", err)
		return
	}
	for _, p := range overdue {
		credit, err := cr.GetByID(ctx, p.CreditID)
		if err != nil {
			continue
		}
		acc, err := ar.GetByID(ctx, credit.AccountID)
		if err != nil {
			continue
		}
		if acc.Balance >= p.Amount {
			tx, _ := db.BeginTx(ctx, nil)
			_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", p.Amount, credit.AccountID)
			if err == nil {
				tx.ExecContext(ctx, "UPDATE payment_schedules SET status='paid' WHERE id=$1", p.ID)
				tx.Commit()
				// Отправка email
				email, err := ur.GetEmailByID(ctx, credit.UserID)
				if err == nil {
					es.SendPaymentNotification(email, p.Amount)
				}
			} else {
				tx.Rollback()
				penalty := p.Amount * 0.1
				newAmount := p.Amount + penalty
				newDue := now.AddDate(0, 1, 0)
				db.ExecContext(ctx, "UPDATE payment_schedules SET amount=$1, due_date=$2, status='overdue' WHERE id=$3", newAmount, newDue, p.ID)
			}
		} else {
			penalty := p.Amount * 0.1
			newAmount := p.Amount + penalty
			newDue := now.AddDate(0, 1, 0)
			db.ExecContext(ctx, "UPDATE payment_schedules SET amount=$1, due_date=$2, status='overdue' WHERE id=$3", newAmount, newDue, p.ID)
		}
	}
}
