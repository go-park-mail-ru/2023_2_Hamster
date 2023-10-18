package user

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/delivery/http/transfer_models"
	"github.com/google/uuid"
)

// Bussiness logic methods to work with user
type Usecase interface {
	//GetByID(userID uuid.UUID) (*models.User, error)
	//	ChangeInfo(user *models.User) error
	GetUserBalance(userID uuid.UUID) (float64, error)
	GetPlannedBudget(userID uuid.UUID) (float64, error)
	GetCurrentBudget(userID uuid.UUID) (float64, error)
	GetAccounts(userID uuid.UUID) ([]models.Accounts, error)
	GetFeed(userID uuid.UUID) (*transfer_models.UserFeed, error)
	GetUser(userID uuid.UUID) (*models.User, error)
}

type Repository interface {
	GetByID(userID uuid.UUID) (*models.User, error)
	CreateUser(user models.User) (uuid.UUID, error)
	//	IncreaseUserVersion(ctx context.Context, userID uuid.UUID) error
	GetUserByUsername(username string) (*models.User, error)
	//	GetUserByIDAndVersion(ctx context.Context, userID, userVersion uuid.UUID) (*models.User, error)
	GetUserBalance(userID uuid.UUID) (float64, error) // transfer account repostiory
	GetPlannedBudget(userID uuid.UUID) (float64, error)
	GetCurrentBudget(userID uuid.UUID) (float64, error)
	GetAccounts(userID uuid.UUID) ([]models.Accounts, error) // transfer account repository
	//IncreaseUserVersion(ctx context.Context, userID uuid.UUID) error
}
