package http

import (
	"html"
	"time"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

const (
	TransactionNotCreate = "can't create transaction"
)

type TransactionCreateResponse struct {
	TransactionID uuid.UUID `json:"transaction_id"`
}

type MasTransaction struct {
	Transactions []models.TransactionTransfer `json:"transaction"`
}

type CreateTransaction struct {
	CategoryID  uuid.UUID `json:"category_id" valid:"required"`
	AccountID   uuid.UUID `json:"account_id" valid:"required"`
	Total       float64   `json:"total" valid:"required,greaterzero"`
	IsIncome    bool      `json:"is_income" valid:"required"`
	Date        time.Time `json:"date" valid:""`
	Payer       string    `json:"payer" valid:""`
	Description string    `json:"description"`
}

func (cr *CreateTransaction) CheckValid() error {
	cr.Payer = html.EscapeString(cr.Payer)
	cr.Description = html.EscapeString(cr.Description)
	_, err := valid.ValidateStruct(*cr)

	return err
}

func (cr *CreateTransaction) ToTransaction(user *models.User) *models.Transaction {
	return &models.Transaction{
		UserID:      user.ID,
		CategoryID:  cr.CategoryID,
		AccountID:   cr.AccountID,
		Total:       cr.Total,
		IsIncome:    cr.IsIncome,
		Date:        cr.Date,
		Description: cr.Description,
	}
}
