package user

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microsevices/user/delivery/http/transfer_models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

// Bussiness logic methods to work with user
type Usecase interface {
	//GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	//	ChangeInfo(user *models.User) error
	GetUserBalance(ctx context.Context, userID uuid.UUID) (float64, error)
	GetPlannedBudget(ctx context.Context, userID uuid.UUID) (float64, error)
	GetCurrentBudget(ctx context.Context, userID uuid.UUID) (float64, error)
	GetAccounts(ctx context.Context, userID uuid.UUID) ([]models.Accounts, error)
	GetFeed(ctx context.Context, userID uuid.UUID) (*transfer_models.UserFeed, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

type Repository interface {
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, user models.User) (uuid.UUID, error)
	//	IncreaseUserVersion(ctx context.Context, ctx context.Context, userID uuid.UUID) error
	GetUserByUsername(username string) (*models.User, error)
	//	GetUserByIDAndVersion(ctx context.Context, ctx context.Context, userID, userVersion uuid.UUID) (*models.User, error)
	GetUserBalance(ctx context.Context, userID uuid.UUID) (float64, error) // transfer account repostiory
	GetPlannedBudget(ctx context.Context, userID uuid.UUID) (float64, error)
	GetCurrentBudget(ctx context.Context, userID uuid.UUID) (float64, error)
	GetAccounts(ctx context.Context, userID uuid.UUID) ([]models.Accounts, error) // transfer account repository
	//IncreaseUserVersion(ctx context.Context, ctx context.Context, userID uuid.UUID) error
	UpdateUser(ctx context.Context, user *models.User) error
	CheckUser(ctx context.Context, userID uuid.UUID) error
}
