package auth

import (
	"html"

	valid "github.com/asaskevich/govalidator"

	"github.com/google/uuid"
)

// Models for auth Request-Response handling

type (
	Session struct {
		ID     uuid.UUID `json:"user_id"`
		Cookie string    `json:"cookie"`
	}

	SignUpInput struct {
		Login          string `json:"login" valid:"required,length(4|20)"`
		Username       string `json:"username" valid:"required,length(4|20)"`
		PlaintPassword string `json:"password" valid:"required,length(4|20)"`
	}

	LoginInput struct {
		Login          string `json:"login" valid:"required,length(4|20)"`
		PlaintPassword string `json:"password" valid:"required,length(4|20)"`
	}

	SignResponse struct {
		ID       uuid.UUID `json:"id" valid:"required"`
		Username string    `json:"username" valid:"required"`
	}
)

func (li *LoginInput) CheckValid() error {
	li.Login = html.EscapeString(li.Login)
	li.PlaintPassword = html.EscapeString(li.PlaintPassword)

	_, err := valid.ValidateStruct(*li)

	return err
}

func (si *SignUpInput) CheckValid() error {
	si.Login = html.EscapeString(si.Login)
	si.Username = html.EscapeString(si.Username)
	si.PlaintPassword = html.EscapeString(si.PlaintPassword)

	_, err := valid.ValidateStruct(*si)
	return err
}
