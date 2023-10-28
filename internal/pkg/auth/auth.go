package auth

import (
	"context"

	"github.com/google/uuid"
)

type Usecase interface {
	// Auth
	SignUp(ctx context.Context, input SignUpInput) (uuid.UUID, string, error)
	Login(ctx context.Context, login, plainPassword string) (uuid.UUID, string, error)
}

type Repository interface {
	// Create User
	// CreateUser(ctx context.Context, user models.User) (models.User, error)

	// Validation
	CheckCorrectPassword(ctx context.Context, password string) error
	CheckExistUsername(ctx context.Context, username string) error
	CheckLoginUnique(ctx context.Context, login string) (bool, error)
}
