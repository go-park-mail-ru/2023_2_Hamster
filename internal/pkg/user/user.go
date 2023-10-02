package user

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

// Bussiness logic methods to work with user
type Usecase interface {
	GetByID(userID uuid.UUID) (*models.User, error)
	ChangeInfo(user *models.User) error
	GetBalance(user *models.User) (float32, error)
}

type Repository interface {
	GetByID(userID uuid.UUID) (*models.User, error)
	CreateUser(user models.User) (uuid.UUID, error)
	GetUserByUsername(username string) (*models.User, error)

	//IncreaseUserVersion(ctx context.Context, userID uuid.UUID) error
}
