package transaction

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Usecase interface {
	DeleteTransaction(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID) error
	CreateTransaction(ctx context.Context, transaction *models.Transaction) (uuid.UUID, error)
	// GetTransaction(ctx context.Context, transaction models.Transaction) *models.Transaction
	GetFeed(ctx context.Context, userID uuid.UUID, page int, pageSize int) ([]models.Transaction, bool, error)
	UpdateTransaction(ctx context.Context, transaction *models.Transaction) error
}

type Repository interface {
	DeleteTransaction(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID) error
	CreateTransaction(ctx context.Context, transaction *models.Transaction) (uuid.UUID, error)
	GetFeed(ctx context.Context, userID uuid.UUID, page int, pageSize int) ([]models.Transaction, bool, error)
	// GetTransaction(ctx context.Context, transaction models.Transaction) *models.Transaction
	UpdateTransaction(ctx context.Context, transaction *models.Transaction) error
	CheckTransaciont(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID) error
}
