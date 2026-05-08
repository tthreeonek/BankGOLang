package postgres

import (
	"bank-api/internal/model"
	"context"
	"database/sql"
)

type AccountRepo struct {
	DB *sql.DB
}

func NewAccountRepo(db *sql.DB) *AccountRepo {
	return &AccountRepo{DB: db}
}

func (r *AccountRepo) Create(ctx context.Context, a *model.Account) error {
	return r.DB.QueryRowContext(ctx,
		"INSERT INTO accounts (user_id, balance) VALUES ($1, $2) RETURNING id, created_at",
		a.UserID, a.Balance,
	).Scan(&a.ID, &a.CreatedAt)
}

func (r *AccountRepo) GetByID(ctx context.Context, id int) (*model.Account, error) {
	a := &model.Account{}
	err := r.DB.QueryRowContext(ctx,
		"SELECT id, user_id, balance, created_at FROM accounts WHERE id = $1", id,
	).Scan(&a.ID, &a.UserID, &a.Balance, &a.CreatedAt)
	return a, err
}

func (r *AccountRepo) GetByUserID(ctx context.Context, userID int) ([]model.Account, error) {
	rows, err := r.DB.QueryContext(ctx,
		"SELECT id, user_id, balance, created_at FROM accounts WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accs []model.Account
	for rows.Next() {
		var a model.Account
		if err := rows.Scan(&a.ID, &a.UserID, &a.Balance, &a.CreatedAt); err != nil {
			return nil, err
		}
		accs = append(accs, a)
	}
	return accs, nil
}

func (r *AccountRepo) UpdateBalance(ctx context.Context, id int, delta float64) error {
	_, err := r.DB.ExecContext(ctx,
		"UPDATE accounts SET balance = balance + $1 WHERE id = $2", delta, id)
	return err
}

func (r *AccountRepo) Transfer(ctx context.Context, fromID, toID int, amount float64) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, fromID)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toID)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		"INSERT INTO transactions (from_account_id, to_account_id, amount, type) VALUES ($1,$2,$3,'transfer')",
		fromID, toID, amount)
	if err != nil {
		return err
	}
	return tx.Commit()
}
