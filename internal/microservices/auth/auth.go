package auth

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

// go:generate mockgen -source=auth.go -destination=mocks/auth_mock.go
// go:generate protoc --go_out=.  --go-grpc_out=.  --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative auth.proto

type Usecase interface {
	// Auth
	SignUp(ctx context.Context, input SignUpInput) (uuid.UUID, string, string, error)
	Login(ctx context.Context, login, plainPassword string) (uuid.UUID, string, string, error)
	CheckLoginUnique(ctx context.Context, login string) (bool, error)

	ChangePassword(ctx context.Context, input ChangePasswordInput) error

	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

type Repository interface {
	// Create User
	CreateUser(ctx context.Context, u models.User) (uuid.UUID, error)

	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)

	ChangePassword(ctx context.Context, userID uuid.UUID, newPassword string) error

	// Validation
	// CheckCorrectPassword(ctx context.Context, password string) error
	// CheckExistUsername(ctx context.Context, username string) error
	CheckLoginUnique(ctx context.Context, login string) (bool, error)
}
