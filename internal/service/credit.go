package service

import (
	"context"
	"math"
	"time"

	"bank-api/internal/model"
	"bank-api/internal/repository"

	"github.com/sirupsen/logrus"
)

type CreditService struct {
	creditRepo repository.CreditRepository
	cbrService *CBRService
}

func NewCreditService(cr repository.CreditRepository, cbr *CBRService) *CreditService {
	return &CreditService{creditRepo: cr, cbrService: cbr}
}

func (s *CreditService) CreateCredit(ctx context.Context, userID, accountID int, amount float64, termMonths int) (*model.Credit, error) {
	rate, err := s.cbrService.GetKeyRate() // ставка ЦБ + маржа
	if err != nil {
		logrus.Error("Failed to get key rate, using default 15%")
		rate = 15.0 // fallback
	}
	monthlyRate := rate / 100.0 / 12.0
	// аннуитетный платёж
	monthlyPayment := amount * monthlyRate * math.Pow(1+monthlyRate, float64(termMonths)) /
		(math.Pow(1+monthlyRate, float64(termMonths)) - 1)

	credit := &model.Credit{
		UserID:         userID,
		AccountID:      accountID,
		Amount:         amount,
		InterestRate:   rate,
		TermMonths:     termMonths,
		MonthlyPayment: math.Round(monthlyPayment*100) / 100,
		StartDate:      time.Now(),
		Status:         "active",
	}
	if err := s.creditRepo.Create(ctx, credit); err != nil {
		return nil, err
	}

	// Генерация графика платежей
	schedule := make([]model.PaymentSchedule, termMonths)
	for i := 0; i < termMonths; i++ {
		dueDate := time.Now().AddDate(0, i+1, 0)
		schedule[i] = model.PaymentSchedule{
			CreditID: credit.ID,
			DueDate:  dueDate,
			Amount:   credit.MonthlyPayment,
			Status:   "pending",
		}
	}
	if err := s.creditRepo.CreatePaymentSchedule(ctx, schedule); err != nil {
		return nil, err
	}
	return credit, nil
}

func (s *CreditService) GetSchedule(ctx context.Context, creditID int) ([]model.PaymentSchedule, error) {
	return s.creditRepo.GetScheduleByCreditID(ctx, creditID)
}
