package models

import "github.com/google/uuid"

type Accounts struct {
	ID          uuid.UUID `json:"id"`
	username    uint      `json:"user_id"`
	Balance     float64   `json:"balance"`
	MeanPayment string    `json:"mean_payment"`
}
