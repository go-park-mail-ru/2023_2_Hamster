package auth

import "github.com/google/uuid"

type (
	Session struct {
		ID     uuid.UUID `json:"user_id"`
		Cookie string    `json:"cookie"`
	}

	signUpResponse struct {
		ID       uuid.UUID `json:"id"`
		Username string    `json:"username"`
	}
)
