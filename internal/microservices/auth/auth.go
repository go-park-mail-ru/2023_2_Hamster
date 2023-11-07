package auth

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Usecase interface {
	// Auth
	SignUp(ctx context.Context, input SignUpInput) (uuid.UUID, string, error)
	Login(ctx context.Context, login, plainPassword string) (uuid.UUID, string, error)
	CheckLoginUnique(ctx context.Context, login string) (bool, error)
}

type Repository interface {
	// Create User
	CreateUser(ctx context.Context, u models.User) (uuid.UUID, error)

	GetUserByLogin(ctx context.Context, login string) (*models.User, error)

	// Validation
	// CheckCorrectPassword(ctx context.Context, password string) error
	CheckExistUsername(ctx context.Context, username string) error
	CheckLoginUnique(ctx context.Context, login string) (bool, error)
}
