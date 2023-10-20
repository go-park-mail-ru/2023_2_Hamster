package http

import "github.com/google/uuid"

type signUpResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type signInput struct {
	Username string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required"`
}

type loginResponse struct { // same response for sign up
	JWT string `json:"access_token"`
}
