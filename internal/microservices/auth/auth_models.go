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

	UserIdInput struct {
		ID uuid.UUID `json:"user_id"`
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
		Login    string    `json:"login" valid:"required"`
		Username string    `json:"username" valid:"required"`
	}

	UniqCheckInput struct {
		Login string `json:"login" valid:"required"`
	}

	ChangePasswordInput struct {
		Login       string
		OldPassword string `json:"old_password" valid:"required,length(4|20)"`
		NewPassword string `json:"new_password" valid:"required,length(4|20)"`
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

func (ui *UniqCheckInput) CheckValid() error {
	ui.Login = html.EscapeString(ui.Login)
	_, err := valid.ValidateStruct(*ui)
	return err
}

func (ci *ChangePasswordInput) CheckValid() error {
	ci.OldPassword = html.EscapeString(ci.OldPassword)
	ci.NewPassword = html.EscapeString(ci.NewPassword)
	_, err := valid.ValidateStruct(*ci)
	return err
}

// Errors

const (
	_                     = iota
	InternalDataBaseError = "internal database error"
	InvalidBodyRequest    = "invalid input params"
	ForbiddenUser         = "user has no rights"
)

// type customErr struct {
// 	Err error
// 	Msg string
// }
//
// func (e *customErr) Error() string {
// 	return e.Msg
// }
//
// func (e *customErr) Unwrap() error {
// 	return e.Err
// }
//
