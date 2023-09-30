package models

import (
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type Transaction struct {
	ID          uuid.UUID `json:"id" valid:"-"`
	UserID      uint      `json:"user_id" valid:"-"`
	CategoryID  uint      `json:"category_id" valid:"-"`
	AccountID   uint      `json:"account_id" valid:"-"`
	Total       float64   `json:"total" valid:"required,greaterzero"`
	IsIncome    bool      `json:"is_income" valid:"required"`
	Date        time.Time `json:"date" valid:"isdate"`
	Payer       string    `json:"payer" valid:"payer"`
	Description string    `json:"description" valid:"-"`
}

func (t *Transaction) TransactionValidate() error {
	_, err := valid.ValidateStruct(t)
	return err
}
