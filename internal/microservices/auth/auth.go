package auth

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Usecase interface {
	// Auth
	SignUp(ctx context.Context, input models.SignInput) (models.User, error)
	Login(ctx context.Context, input models.SignInput) (models.User, error)

	// Session
	GetSessionByCookie(ctx context.Context, cookie string) (Session, error)
	CreateSessionById(ctx context.Context, userID uuid.UUID) (Session, error)
	DeleteSessionByCookie(ctx context.Context, cookie string) error
}

type Repository interface {
	// Create User
	CreateUser(ctx context.Context, user models.User) (models.User, error)

	// Validation
	CheckCorrectPassword(ctx context.Context, password string) error
	CheckExistUsername(ctx context.Context, username string) error

	// Session
	GetSessionByCookie(ctx context.Context, cookie string) (Session, error)
	CreateSession(ctx context.Context, session Session) error
	DeleteSession(ctx context.Context, cookie string) error
}
