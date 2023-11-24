package models

import "github.com/google/uuid"

type Accounts struct {
	ID             uuid.UUID `json:"id"`
	Balance        float64   `json:"balance"`
	Accumulation   bool      `json:"accumulation"`
	BalanceEnabled bool      `json:"balance_enabled"`
	MeanPayment    string    `json:"mean_payment"`
}

type AccounstTransfer struct {
	ID             uuid.UUID `json:"id"`
	Balance        float64   `json:"balance"`
	Accumulation   bool      `json:"accumulation"`
	BalanceEnabled bool      `json:"balance_enabled"`
	MeanPayment    string    `json:"mean_payment"`
}
