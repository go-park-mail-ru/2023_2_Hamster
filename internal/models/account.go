package models

import "github.com/google/uuid"

type Accounts struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uint      `json:"user_id" db:"user_id"`
	Balance     float64   `json:"balance" db:"balance"`
	MeanPayment string    `json:"mean_payment" db:"mean_payment"`
}
