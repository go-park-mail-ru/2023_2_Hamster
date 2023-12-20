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

type TransactionCount struct {
	Count int `json:"count"`
}

type MasTransaction struct {
	Transactions []models.TransactionTransfer `json:"transactions"`
}

//easyjson:json
type CreateTransaction struct {
	AccountIncomeID  uuid.UUID             `json:"account_income" valid:"-"`  // ???
	AccountOutcomeID uuid.UUID             `json:"account_outcome" valid:"-"` // ???
	Income           float64               `json:"income" valid:"-"`
	Outcome          float64               `json:"outcome" valid:"-"`
	Date             time.Time             `json:"date" valid:"required"`
	Payer            string                `json:"payer," valid:"maxstringlength(20)"`
	Description      string                `json:"description,omitempty" valid:""`
	Categories       []models.CategoryName `json:"categories" valid:"-"`
}

//easyjson:json
type UpdTransaction struct {
	ID               uuid.UUID             `json:"transaction_id" valid:"required"`
	AccountIncomeID  uuid.UUID             `json:"account_income" valid:"-"`
	AccountOutcomeID uuid.UUID             `json:"account_outcome" valid:"-"`
	Income           float64               `json:"income" valid:"-"`
	Outcome          float64               `json:"outcome" valid:"-"`
	Date             time.Time             `json:"date" valid:"required"`
	Payer            string                `json:"payer" valid:"maxstringlength(20)"`
	Description      string                `json:"description" valid:"-"`
	Categories       []models.CategoryName `json:"categories"`
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

func (cr *CreateTransaction) ToTransaction(user *models.User) *models.Transaction {
	return &models.Transaction{
		UserID:           user.ID,
		AccountIncomeID:  cr.AccountIncomeID,
		AccountOutcomeID: cr.AccountOutcomeID,
		Income:           cr.Income,
		Outcome:          cr.Outcome,
		Payer:            cr.Payer,
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
		Payer:            ut.Payer,
		Outcome:          ut.Outcome,
		Date:             ut.Date,
		Description:      ut.Description,
		Categories:       ut.Categories,
	}
}
