package models

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id" valid:"-"`
	Username      string    `json:"name" valid:"-"`
	FirstName     string    `json:"first_name" valid:"required,runelength(2|20)"`
	LastName      string    `json:"last_name" valid:"required,runelength(2|20)"`
	PlannedBudget float64   `json:"planned_budget" valid:"-"`
	Password      string    `json:"password" valid:"required,runelength(7|30),passwordcheck"`
	AvatarURL     string    `json:"avatar_url" vaild:"-"`
	Salt          string    `json:"salt"`
}

func (u *User) UserValidate() error {
	_, err := valid.ValidateStruct(u)
	return err
}
