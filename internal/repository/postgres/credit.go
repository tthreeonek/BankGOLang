package postgres

import (
	"context"
	"database/sql"
	"time"

	"bank-api/internal/model"
)

type CreditRepo struct {
	DB *sql.DB
}

func NewCreditRepo(db *sql.DB) *CreditRepo {
	return &CreditRepo{DB: db}
}

func (r *CreditRepo) Create(ctx context.Context, c *model.Credit) error {
	return r.DB.QueryRowContext(ctx,
		`INSERT INTO credits (user_id, account_id, amount, interest_rate, term_months, monthly_payment, start_date, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		c.UserID, c.AccountID, c.Amount, c.InterestRate, c.TermMonths, c.MonthlyPayment, c.StartDate, c.Status,
	).Scan(&c.ID)
}

func (r *CreditRepo) GetByID(ctx context.Context, id int) (*model.Credit, error) {
	c := &model.Credit{}
	err := r.DB.QueryRowContext(ctx,
		"SELECT id, user_id, account_id, amount, interest_rate, term_months, monthly_payment, start_date, status FROM credits WHERE id=$1",
		id).Scan(&c.ID, &c.UserID, &c.AccountID, &c.Amount, &c.InterestRate, &c.TermMonths, &c.MonthlyPayment, &c.StartDate, &c.Status)
	return c, err
}

func (r *CreditRepo) CreatePaymentSchedule(ctx context.Context, schedule []model.PaymentSchedule) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, p := range schedule {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO payment_schedules (credit_id, due_date, amount, status) VALUES ($1,$2,$3,$4)",
			p.CreditID, p.DueDate, p.Amount, p.Status)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *CreditRepo) GetOverduePayments(ctx context.Context, now time.Time) ([]model.PaymentSchedule, error) {
	rows, err := r.DB.QueryContext(ctx,
		"SELECT id, credit_id, due_date, amount, status FROM payment_schedules WHERE due_date <= $1 AND status = 'pending'",
		now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PaymentSchedule
	for rows.Next() {
		var p model.PaymentSchedule
		if err := rows.Scan(&p.ID, &p.CreditID, &p.DueDate, &p.Amount, &p.Status); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *CreditRepo) MarkPaymentPaid(ctx context.Context, id int) error {
	_, err := r.DB.ExecContext(ctx, "UPDATE payment_schedules SET status = 'paid' WHERE id = $1", id)
	return err
}

// GetScheduleByCreditID возвращает все платежи по кредиту, отсортированные по дате
func (r *CreditRepo) GetScheduleByCreditID(ctx context.Context, creditID int) ([]model.PaymentSchedule, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, credit_id, due_date, amount, status
		 FROM payment_schedules
		 WHERE credit_id = $1
		 ORDER BY due_date`, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedule []model.PaymentSchedule
	for rows.Next() {
		var p model.PaymentSchedule
		if err := rows.Scan(&p.ID, &p.CreditID, &p.DueDate, &p.Amount, &p.Status); err != nil {
			return nil, err
		}
		schedule = append(schedule, p)
	}
	return schedule, rows.Err()
}
