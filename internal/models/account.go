package models

import "github.com/google/uuid"

type Accounts struct {
	ID             uuid.UUID     `json:"id"`
	Balance        float64       `json:"balance"`
	Accumulation   bool          `json:"accumulation"`
	SharingID      uuid.UUID     `json:"sharing_id"`
	BalanceEnabled bool          `json:"balance_enabled"`
	MeanPayment    string        `json:"mean_payment"`
	Users          []SharingUser `json:"users"`
}

type AccounstTransfer struct {
	ID             uuid.UUID `json:"id"`
	Balance        float64   `json:"balance"`
	Accumulation   bool      `json:"accumulation"`
	BalanceEnabled bool      `json:"balance_enabled"`
	MeanPayment    string    `json:"mean_payment"`
}

type AddUserAccount struct {
	Login     string    `json:"login"`
	AccountID uuid.UUID `json:"account_id"`
}

type DeleteInAccount struct {
	UserID    uuid.UUID `json:"user_id"`
	AccountID uuid.UUID `json:"account_id"`
}
