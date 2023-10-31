package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID               uuid.UUID   `json:"id" valid:"-"`
	UserID           uuid.UUID   `json:"user_id" valid:"-"`
	AccountIncomeID  uuid.UUID   `json:"account_income" valid:"-"`
	AccountOutcomeID uuid.UUID   `json:"account_outcome" valid:"-"`
	Income           float64     `json:"income" valid:"required"`
	Outcome          float64     `json:"outcome" valid:"required"`
	Date             time.Time   `json:"date" valid:"isdate"`
	Payer            string      `json:"payer" valid:"payer"`
	Description      string      `json:"description" valid:"-"`
	Categories       []uuid.UUID `json:"categories" valid:"-"`
}

type TransactionTransfer struct {
	ID               uuid.UUID   `json:"id" valid:"-"`
	AccountIncomeID  uuid.UUID   `json:"account_income" valid:"-"`
	AccountOutcomeID uuid.UUID   `json:"account_outcome" valid:"-"`
	Income           float64     `json:"income" valid:"required"`
	Outcome          float64     `json:"outcome" valid:"required"`
	Date             time.Time   `json:"date" valid:"isdate"`
	Payer            string      `json:"payer" valid:"payer"`
	Description      string      `json:"description" valid:"-"`
	Categories       []uuid.UUID `json:"categories" valid:"-"`
}

func InitTransactionTransfer(transaction Transaction) TransactionTransfer {
	return TransactionTransfer{
		ID:               transaction.ID,
		AccountIncomeID:  transaction.AccountIncomeID,
		AccountOutcomeID: transaction.AccountOutcomeID,
		Income:           transaction.Income,
		Outcome:          transaction.Outcome,
		Date:             transaction.Date,
		Payer:            transaction.Payer,
		Description:      transaction.Description,
		Categories:       transaction.Categories,
	}
}
