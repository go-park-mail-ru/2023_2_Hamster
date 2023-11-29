package account

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Usecase interface {
	CreateAccount(ctx context.Context, userID uuid.UUID, account *models.Accounts) (uuid.UUID, error)
	UpdateAccount(ctx context.Context, userID uuid.UUID, account *models.Accounts) error
	DeleteAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error
}

type Repository interface {
	CreateAccount(ctx context.Context, userID uuid.UUID, account *models.Accounts) (uuid.UUID, error)
	UpdateAccount(ctx context.Context, userID uuid.UUID, account *models.Accounts) error
	DeleteAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error
	CheckForbidden(ctx context.Context, accountID uuid.UUID, userID uuid.UUID) error
}
