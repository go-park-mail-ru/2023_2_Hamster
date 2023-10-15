package models

import "github.com/google/uuid"

type Accounts struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Balance     float64   `json:"balance"`
	MeanPayment string    `json:"mean_payment"`
}

type AccounstTransfer struct {
	ID            uuid.UUID `json:"id"`
	Username      string    `json:"username"`
	PlannedBudget float64   `json:"planned_budget"`
}
