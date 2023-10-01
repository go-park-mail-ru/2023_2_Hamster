package user

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
)

// Bussiness logic methods to work with user
type Usecase interface {
	GetByID(userID uint32) (*models.User, error)
	ChangeInfo(user *models.User) error
}

type Repository interface {
	GetByID(userID uint32) (*models.User, error)
	CreateUser(user models.User) (uint32, error)

	IncreaseUserVersion(ctx context.Context, userID uint32) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByIDAndVersion(ctx context.Context, userID, userVersion uint32) (*models.User, error)
}
