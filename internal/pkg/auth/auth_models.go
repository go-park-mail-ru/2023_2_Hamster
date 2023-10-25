package auth

import "github.com/google/uuid"

type (
	Session struct {
		ID     uuid.UUID `json:"user_id"`
		Cookie string    `json:"cookie"`
	}

	SignResponse struct {
		ID       uuid.UUID `json:"user_id"`
		Username string    `json:"username"`
	}

	SignUser struct {
		Username       string `json:"username"`
		PlaintPassword string `json:"password"`
	}
)
