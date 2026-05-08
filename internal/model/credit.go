package model

import "time"

type Credit struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	AccountID      int       `json:"account_id"`
	Amount         float64   `json:"amount"`
	InterestRate   float64   `json:"interest_rate"`
	TermMonths     int       `json:"term_months"`
	MonthlyPayment float64   `json:"monthly_payment"`
	StartDate      time.Time `json:"start_date"`
	Status         string    `json:"status"` // active, closed, overdue
}

type PaymentSchedule struct {
	ID       int       `json:"id"`
	CreditID int       `json:"credit_id"`
	DueDate  time.Time `json:"due_date"`
	Amount   float64   `json:"amount"`
	Status   string    `json:"status"` // pending, paid, overdue
}
