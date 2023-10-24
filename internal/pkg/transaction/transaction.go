package transactoin

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Usecase interface {
	DeleteTransaction(ctx context.Context, transactionID uuid.UUID) error
	CreateTransaction(ctx context.Context, transaction models.Transaction) (uuid.UUID, error)
	GetFeed(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *models.Transaction) error
}

type Repository interface {
	DeleteTransaction(ctx context.Context, transactionID uuid.UUID) error
	CreateTransaction(ctx context.Context, transaction models.Transaction) (uuid.UUID, error)
	GetFeed(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *models.Transaction) error
}
