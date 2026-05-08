package repository

import (
	"context"
	"time"

	"bank-api/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	GetEmailByID(ctx context.Context, id int) (string, error) // <-- добавлено
}

type AccountRepository interface {
	Create(ctx context.Context, account *model.Account) error
	GetByID(ctx context.Context, id int) (*model.Account, error)
	GetByUserID(ctx context.Context, userID int) ([]model.Account, error)
	UpdateBalance(ctx context.Context, id int, delta float64) error
	Transfer(ctx context.Context, fromID, toID int, amount float64) error
}

type CardRepository interface {
	Create(ctx context.Context, card *model.Card) error
	GetByAccountID(ctx context.Context, accountID int) ([]model.Card, error)
}

type CreditRepository interface {
	Create(ctx context.Context, credit *model.Credit) error
	GetByID(ctx context.Context, id int) (*model.Credit, error)
	CreatePaymentSchedule(ctx context.Context, schedule []model.PaymentSchedule) error
	GetOverduePayments(ctx context.Context, now time.Time) ([]model.PaymentSchedule, error)
	MarkPaymentPaid(ctx context.Context, id int) error
	GetScheduleByCreditID(ctx context.Context, creditID int) ([]model.PaymentSchedule, error)
}
