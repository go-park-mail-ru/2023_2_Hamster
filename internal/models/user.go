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

type SignInput struct {
	Username string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required"`
}

type ContextKeyUserType struct{}

func (u *User) UserValidate() error {
	_, err := valid.ValidateStruct(u)
	return err
}

type UserTransfer struct {
	ID            uuid.UUID `json:"id"`
	Username      string    `json:"username"`
	PlannedBudget float64   `json:"planned_budget"`
	AvatarURL     string    `json:"avatar_url"`
}

func InitUserTransfer(user User) UserTransfer {
	return UserTransfer{
		ID:            user.ID,
		Username:      user.Username,
		PlannedBudget: user.PlannedBudget,
		AvatarURL:     user.AvatarURL,
	}
}
