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
	TransactionNotSuch   = "can't such transactoin"

	TransactionCreateServerError = "can't get transaction"
)

type TransactionCreateResponse struct {
	TransactionID uuid.UUID `json:"transaction_id"`
}

type MasTransaction struct {
	Transactions []models.TransactionTransfer `json:"transaction"`
	IsAll        bool                         `json:"is_all"`
}

type CreateTransaction struct {
	CategoryID  uuid.UUID `json:"category_id" valid:"required"`
	AccountID   uuid.UUID `json:"account_id" valid:"required"`
	Total       float64   `json:"total" valid:"required"`
	IsIncome    bool      `json:"is_income" valid:""`
	Date        time.Time `json:"date" valid:"required"`
	Payer       string    `json:"payer" valid:""`
	Description string    `json:"description" valid:""`
}

type UpdTransaction struct {
	ID          uuid.UUID `json:"transaction_id" valid:"required"`
	CategoryID  uuid.UUID `json:"category_id" valid:"required"`
	AccountID   uuid.UUID `json:"account_id" valid:"required"`
	Total       float64   `json:"total" valid:"required"`
	IsIncome    bool      `json:"is_income" valid:""`
	Date        time.Time `json:"date" valid:"required"`
	Payer       string    `json:"payer" valid:"required"`
	Description string    `json:"description" valid:"required"`
}

func (cr *CreateTransaction) CheckValid() error {
	cr.Payer = html.EscapeString(cr.Payer)
	cr.Description = html.EscapeString(cr.Description)

	_, err := valid.ValidateStruct(*cr)

	return err
}

func (cr *UpdTransaction) CheckValid() error {
	cr.Payer = html.EscapeString(cr.Payer)
	cr.Description = html.EscapeString(cr.Description)

	_, err := valid.ValidateStruct(*cr)

	return err
}

type QueryListOptions struct {
	Page     int `json:"page" minimum:"1" validate:"optional" example:"1"`
	PageSize int `json:"page_size" minimum:"1" maximum:"20" validate:"optional" example:"10"`
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

func (ut *UpdTransaction) ToTransaction(user *models.User) *models.Transaction {
	return &models.Transaction{
		ID:          ut.ID,
		UserID:      user.ID,
		CategoryID:  ut.CategoryID,
		AccountID:   ut.AccountID,
		Total:       ut.Total,
		IsIncome:    ut.IsIncome,
		Date:        ut.Date,
		Description: ut.Description,
	}
}
