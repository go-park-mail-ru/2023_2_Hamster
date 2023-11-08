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
	TransactionDeleteServerError = "cat't delete transaction"
)

type TransactionCreateResponse struct {
	TransactionID uuid.UUID `json:"transaction_id"`
}

type MasTransaction struct {
	Transactions []models.TransactionTransfer `json:"transactions"`
	IsAll        bool                         `json:"is_all"`
}

type CreateTransaction struct {
	AccountIncomeID  uuid.UUID   `json:"account_income" valid:"required"`
	AccountOutcomeID uuid.UUID   `json:"account_outcome" valid:"required"`
	Income           float64     `json:"income" valid:"-"`
	Outcome          float64     `json:"outcome" valid:"-"`
	Date             time.Time   `json:"date" valid:"required"`
	Payer            string      `json:"payer" valid:"maxstringlength(20)"`
	Description      string      `json:"description" valid:""`
	Categories       []uuid.UUID `json:"categories" valid:"-"`
}

type UpdTransaction struct {
	ID               uuid.UUID   `json:"transaction_id" valid:"required"`
	AccountIncomeID  uuid.UUID   `json:"account_income" valid:"-"`
	AccountOutcomeID uuid.UUID   `json:"account_outcome" valid:"-"`
	Income           float64     `json:"income" valid:"-"`
	Outcome          float64     `json:"outcome" valid:"-"`
	Date             time.Time   `json:"date" valid:"required"`
	Payer            string      `json:"payer" valid:"maxstringlength(20)"`
	Description      string      `json:"description" valid:"-"`
	Categories       []uuid.UUID `json:"categories"`
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
		UserID:           user.ID,
		AccountIncomeID:  cr.AccountIncomeID,
		AccountOutcomeID: cr.AccountOutcomeID,
		Income:           cr.Income,
		Outcome:          cr.Outcome,
		Date:             cr.Date,
		Description:      cr.Description,
		Categories:       cr.Categories,
	}
}

func (ut *UpdTransaction) ToTransaction(user *models.User) *models.Transaction {
	return &models.Transaction{
		ID:               ut.ID,
		UserID:           user.ID,
		AccountIncomeID:  ut.AccountIncomeID,
		AccountOutcomeID: ut.AccountOutcomeID,
		Income:           ut.Income,
		Outcome:          ut.Outcome,
		Date:             ut.Date,
		Description:      ut.Description,
		Categories:       ut.Categories,
	}
}
