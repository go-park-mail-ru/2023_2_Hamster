package http

import (
	"html"

	valid "github.com/asaskevich/govalidator"
)

type CreateAccount struct {
	Balance        float64 `json:"balance" valid:"required"`
	Accumulation   bool    `json:"accumulation" valid:"required"`
	BalanceEnabled bool    `json:"balance_enabled" valid:"required"`
	MeanPayment    string  `json:"mean_payment" valid:"required,length(1|30)"`
}

func (ca *CreateAccount) CheckValid() error {
	ca.MeanPayment = html.EscapeString(ca.MeanPayment)

	_, err := valid.ValidateStruct(*ca)

	return err
}
