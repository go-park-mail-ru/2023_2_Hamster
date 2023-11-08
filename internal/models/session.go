package models

import "github.com/google/uuid"

type Session struct {
	UserId uuid.UUID `json:"user_id" db:"profile_id"`
	Cookie string    `json:"cookie" db:"cookie"`
}
