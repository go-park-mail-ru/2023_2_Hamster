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
	Payer            string      `json:"payer" valid:"-"`
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
	Payer            string      `json:"payer" valid:"-"`
	Description      string      `json:"description" valid:"-"`
	Categories       []uuid.UUID `json:"categories" valid:"-"`
}

type QueryListOptions struct {
	Category  uuid.UUID `json:"category" validate:"optional" example:"uuid"`
	Account   uuid.UUID `json:"account" validate:"optional" example:"uuid"`
	Income    bool      `json:"income" validate:"optional" example:"true"`
	Outcome   bool      `json:"outcome" validate:"optional" example:"true"`
	StartDate time.Time `json:"start_date" validate:"optional" example:"2023-11-21T19:30:57+03:00"`
	EndDate   time.Time `json:"end_date" validate:"optional" example:"2023-12-21T19:30:57+03:00"`
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
