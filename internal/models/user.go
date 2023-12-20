package models

import (
	"fmt"

	"github.com/google/uuid"
)

// models users
// model ls
type User struct {
	ID            uuid.UUID `json:"id"`
	Login         string    `json:"login"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	PlannedBudget float64   `json:"planned_budget"`
	AvatarURL     uuid.UUID `json:"avatar_url"`
}

type SharingUser struct {
	ID        uuid.UUID `json:"id"`
	Login     string    `json:"login"`
	AvatarURL uuid.UUID `json:"avatar_url"`
}

type ContextKeyUserType struct{}

type UserAlreadyExistsError struct{}

func (e *UserAlreadyExistsError) Error() string {
	return "user already exists"
}

type IncorrectPasswordError struct {
	UserID uuid.UUID
}

func (e *IncorrectPasswordError) Error() string {
	if e.UserID == uuid.Nil {
		return "incorrect password for user"
	}
	return fmt.Sprintf("incorrect password for user #%d", e.UserID)
}
