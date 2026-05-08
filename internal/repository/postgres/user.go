package postgres

import (
	"context"
	"database/sql"

	"bank-api/internal/model"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) Create(ctx context.Context, u *model.User) error {
	return r.DB.QueryRowContext(ctx,
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
		u.Username, u.Email, u.PasswordHash,
	).Scan(&u.ID)
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	u := &model.User{}
	err := r.DB.QueryRowContext(ctx,
		"SELECT id, username, email, password_hash FROM users WHERE email = $1", email,
	).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash)
	return u, err
}

func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	return exists, err
}

func (r *UserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", username).Scan(&exists)
	return exists, err
}

func (r *UserRepo) GetEmailByID(ctx context.Context, id int) (string, error) {
	var email string
	err := r.DB.QueryRowContext(ctx, "SELECT email FROM users WHERE id = $1", id).Scan(&email)
	return email, err
}
