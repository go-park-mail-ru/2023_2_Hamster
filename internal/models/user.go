package models

import (
	"github.com/google/uuid"
)

// models users
type User struct {
	ID            uuid.UUID `json:"id"`
	Login         string    `json:"login"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	PlannedBudget float64   `json:"planned_budget"`
	AvatarURL     uuid.UUID `json:"avatar_url"`
}

type ContextKeyUserType struct{}
