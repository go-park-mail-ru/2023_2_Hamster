package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          uuid.UUID `json:"id" valid:"-"`
	UserID      uuid.UUID `json:"user_id" valid:"-"`
	CategoryID  uuid.UUID `json:"category_id" valid:"-"`
	AccountID   uuid.UUID `json:"account_id" valid:"-"`
	Total       float64   `json:"total" valid:"required,greaterzero"`
	IsIncome    bool      `json:"is_income" valid:"required"`
	Date        time.Time `json:"date" valid:"isdate"`
	Payer       string    `json:"payer" valid:"payer"`
	Description string    `json:"description" valid:"-"`
}

type TransactionTransfer struct {
	ID          uuid.UUID `json:"id" valid:"-"`
	CategoryID  uuid.UUID `json:"category_id" valid:"-"`
	AccountID   uuid.UUID `json:"account_id" valid:"-"`
	Total       float64   `json:"total" valid:"required,greaterzero"`
	IsIncome    bool      `json:"is_income" valid:"required"`
	Date        time.Time `json:"date" valid:"isdate"`
	Payer       string    `json:"payer" valid:"payer"`
	Description string    `json:"description" valid:"-"`
}

func InitTransactionTransfer(transaction Transaction) TransactionTransfer {
	return TransactionTransfer{
		ID:          transaction.ID,
		CategoryID:  transaction.CategoryID,
		AccountID:   transaction.AccountID,
		Total:       transaction.Total,
		IsIncome:    transaction.IsIncome,
		Date:        transaction.Date,
		Payer:       transaction.Payer,
		Description: transaction.Description,
	}
}

// func (t *Transaction) TransactionValidate() error {
// 	_, err := valid.ValidateStruct(t)
// 	return err
// }
