package auth

import "github.com/google/uuid"

// Models for auth Request-Response handling

type (
	Session struct {
		ID     uuid.UUID `json:"user_id"`
		Cookie string    `json:"cookie"`
	}

	SignUpInput struct {
		Login          string `json:"login"`
		Username       string `json:"username"`
		PlaintPassword string `json:"password"`
	}

	LoginInput struct {
		Login          string `json:"login"`
		PlaintPassword string `json:"password"`
	}

	SignResponse struct {
		ID       uuid.UUID `json:"user_id"`
		Username string    `json:"username"`
	}
)
