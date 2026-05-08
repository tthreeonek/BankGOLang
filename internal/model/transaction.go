package model

import "time"

type Transaction struct {
	ID        int       `json:"id"`
	From      int       `json:"from_account_id"`
	To        int       `json:"to_account_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"` // transfer, deposit, credit_payment
	CreatedAt time.Time `json:"created_at"`
}
