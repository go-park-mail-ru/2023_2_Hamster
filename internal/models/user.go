package models

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id" valid:"-"`
	Username      string    `json:"username" valid:"-"`
	PlannedBudget float64   `json:"planned_budget" valid:"-"`
	Password      string    `json:"password" valid:"required,runelength(7|30),passwordcheck"`
	AvatarURL     string    `json:"avatar_url" vaild:"-"`
	Salt          string    `json:"salt"`
}

type ContextKeyUserType struct{}

func (u *User) UserValidate() error {
	_, err := valid.ValidateStruct(u)
	return err
}
