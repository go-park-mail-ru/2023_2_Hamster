package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	Login         string    `json:"login"`
	Username      string    `json:"username"`
	PlannedBudget float64   `json:"planned_budget"`
	Password      string    `json:"password"`
	AvatarURL     uuid.UUID `json:"avatar_url"`
	Salt          string    `json:"salt"`
}

type ContextKeyUserType struct{}
