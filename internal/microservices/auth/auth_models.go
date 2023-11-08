package auth

import "github.com/google/uuid"

// Models for auth Request-Response handling

type (
	Session struct {
		ID     uuid.UUID `json:"user_id"`
		Cookie string    `json:"cookie"`
	}

	SignUpInput struct {
		Login          string `json:"login" valid:"required"`
		Username       string `json:"username" valid:"required"`
		PlaintPassword string `json:"password" valid:"required"`
	}

	LoginInput struct {
		Login          string `json:"login" valid:"required"`
		PlaintPassword string `json:"password" valid:"required"`
	}

	SignResponse struct {
		ID       uuid.UUID `json:"id" valid:"required"`
		Username string    `json:"username" valid:"required"`
	}
)
