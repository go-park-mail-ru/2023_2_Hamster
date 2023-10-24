package auth

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
)

type Usecase interface {
	// Auth
	SignUp(ctx context.Context, input models.User) (models.User, error)
	Login(ctx context.Context, input models.SignInput) (models.User, error)
}

type Repository interface {
	// Create User
	CreateUser(ctx context.Context, user models.User) (models.User, error)

	// Validation
	CheckCorrectPassword(ctx context.Context, password string) error
	CheckExistUsername(ctx context.Context, username string) error
}
