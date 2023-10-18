package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	Username      string    `json:"username"`
	PlannedBudget float64   `json:"planned_budget"`
	Password      string    `json:"password"`
	AvatarURL     string    `json:"avatar_url"`
	Salt          string    `json:"salt"`
}

type ContextKeyUserType struct{}

// func (u *User) UserValidate() error {
// 	_, err := valid.ValidateStruct(u)
// 	return err
// }
