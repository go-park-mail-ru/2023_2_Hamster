package auth

import (
	"context"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound       = errors.New("Error user not found")
	ErrInvalidAccessToken = errors.New("Error invalid access token")
)

type CookieToken struct {
	Value   string
	Expires time.Time
}

type Usecase interface {
	// SignUpUser creates new User and returns it's id
	SignUpUser(user models.User) (uuid.UUID, CookieToken, error)

	SignInUser(username, plainPassword string) (uuid.UUID, CookieToken, error)

	// GetUserByCreds returns User if such exist in repository
	GetUserByCreds(ctx context.Context, username, plainPassword string) (*models.User, error)

	// GetUserByAuthData returns User if such exist in repository
	GetUserByAuthData(ctx context.Context, userID uuid.UUID) (*models.User, error)

	GenerateAccessToken(ctx context.Context, user models.User) (CookieToken, error)

	ValidateAccessToken(accessToken string) (uuid.UUID, error)

	// ChangePassword(ctx context.Context, userID uint32, password string) error
}

type Repository interface {
	GetUserByAuthData(ctx context.Context, userID uuid.UUID) (*models.User, error)
	// IncreaseUserVersion(ctx context.Context, userID uint32) error
	// UpdatePassword(ctx context.Context, userID uint32, passwordHash, salt string) error
}
