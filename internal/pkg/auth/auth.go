package auth

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
)

var (
	ErrUserNotFound       = errors.New("Error user not found")
	ErrInvalidAccessToken = errors.New("Error invalid access token")
)

const CtxUserKey = "user"

type Usecase interface {
	// SignUpUser creates new User and returns it's id
	SignUpUser(ctx context.Context, user models.User) (uint32, error)

	// GetUserByCreds returns User if such exist in repository
	GetUserByCreds(ctx context.Context, username, plainPassword string) (*models.User, error)

	LoginUser(username, plainPassword string) (string, error)

	// GetUserByAuthData returns User if such exist in repository
	GetUserByAuthData(ctx context.Context, userID, userVersion uint32) (*models.User, error)

	GenerateAccessToken(ctx context.Context, userID, userVersion uint32) (string, error)

	ValidateAccessToken(accessToken string) (uint32, uint32, error)

	// IncraseUserVersion inc User access token version
	IncreaseUserVersion(ctx context.Context, userID uint32) error

	// ChangePassword(ctx context.Context, userID uint32, password string) error
}

type Repository interface {
	GetUserByAuthData(ctx context.Context, userID, userVaersion uint32) (*models.User, error)
	IncreaseUserVersion(ctx context.Context, userID uint32) error
	UpdatePassword(ctx context.Context, userID uint32, passwordHash, salt string) error
}

type Tabels interface {
	Users() string
}
