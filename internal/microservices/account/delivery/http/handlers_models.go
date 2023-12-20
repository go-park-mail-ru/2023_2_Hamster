package http

import (
	"html"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

const (
	AccountNotCreate = "can't create account"
	AccountNotSuch   = "can't such account"

	AccountCreateServerError     = "can't get account"
	TransactionDeleteServerError = "cat't delete transaction"
)

type AccountCreateResponse struct {
	AccountID uuid.UUID `json:"account_id"`
}

//easyjson:json
type CreateAccount struct {
	Balance        float64 `json:"balance" valid:"-"`
	Accumulation   bool    `json:"accumulation" valid:"-"`
	BalanceEnabled bool    `json:"balance_enabled" valid:"-"`
	MeanPayment    string  `json:"mean_payment" valid:"required,length(1|30)"`
}

//easyjson:json
type UpdateAccount struct {
	ID             uuid.UUID `json:"id" valid:"required"`
	Balance        float64   `json:"balance" valid:""`
	Accumulation   bool      `json:"accumulation" valid:""`
	BalanceEnabled bool      `json:"balance_enabled" valid:""`
	MeanPayment    string    `json:"mean_payment" valid:""`
}

func (cr *CreateAccount) ToAccount() *models.Accounts {
	return &models.Accounts{
		Balance:        cr.Balance,
		Accumulation:   cr.Accumulation,
		BalanceEnabled: cr.BalanceEnabled,
		MeanPayment:    cr.MeanPayment,
	}
}

func (au *UpdateAccount) ToAccount() *models.Accounts {
	return &models.Accounts{
		ID:             au.ID,
		Balance:        au.Balance,
		Accumulation:   au.Accumulation,
		BalanceEnabled: au.BalanceEnabled,
		MeanPayment:    au.MeanPayment,
	}
}

func (ca *CreateAccount) CheckValid() error {
	ca.MeanPayment = html.EscapeString(ca.MeanPayment)

	_, err := valid.ValidateStruct(*ca)

	return err
}

func (ca *UpdateAccount) CheckValid() error {

	_, err := valid.ValidateStruct(*ca)

	return err
}
