package postgres

import (
	"bank-api/internal/model"
	"context"
	"database/sql"
)

type CardRepo struct {
	DB *sql.DB
}

func NewCardRepo(db *sql.DB) *CardRepo {
	return &CardRepo{DB: db}
}

func (r *CardRepo) Create(ctx context.Context, c *model.Card) error {
	return r.DB.QueryRowContext(ctx,
		"INSERT INTO cards (account_id, encrypted_pan, pan_hmac, encrypted_exp, cvv_hash, masked_pan) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id",
		c.AccountID, c.EncryptedPAN, c.PANHMAC, c.EncryptedExp, c.CVVHash, c.MaskedPAN,
	).Scan(&c.ID)
}

func (r *CardRepo) GetByAccountID(ctx context.Context, accountID int) ([]model.Card, error) {
	rows, err := r.DB.QueryContext(ctx,
		"SELECT id, account_id, encrypted_pan, pan_hmac, encrypted_exp, cvv_hash, masked_pan FROM cards WHERE account_id = $1",
		accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cards []model.Card
	for rows.Next() {
		var c model.Card
		if err := rows.Scan(&c.ID, &c.AccountID, &c.EncryptedPAN, &c.PANHMAC, &c.EncryptedExp, &c.CVVHash, &c.MaskedPAN); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, nil
}
