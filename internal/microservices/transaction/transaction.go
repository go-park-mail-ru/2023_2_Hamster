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
	GetFeed(ctx context.Context, userID uuid.UUID, query *models.QueryListOptions) ([]models.Transaction, error)
	GetCount(ctx context.Context, userID uuid.UUID) (int, error)
	UpdateTransaction(ctx context.Context, transaction *models.Transaction) error

	GetTransactionForExport(ctx context.Context, userId uuid.UUID, query *models.QueryListOptions) ([]models.TransactionExport, error)
	// GetTransactionForExport(r.Context(), user.ID, query)
}

type Repository interface {
	DeleteTransaction(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID) error
	CreateTransaction(ctx context.Context, transaction *models.Transaction) (uuid.UUID, error)
	GetFeed(ctx context.Context, userID uuid.UUID, query *models.QueryListOptions) ([]models.Transaction, error)
	GetCount(ctx context.Context, userID uuid.UUID) (int, error)
	// GetTransaction(ctx context.Context, transaction models.Transaction) *models.Transaction
	UpdateTransaction(ctx context.Context, transaction *models.Transaction) error
	CheckForbidden(ctx context.Context, transactinID uuid.UUID) (uuid.UUID, error)
	//Check(ctx context.Context, transactionID uuid.UUID) error

	GetTransactionForExport(ctx context.Context, userId uuid.UUID, query *models.QueryListOptions) ([]models.TransactionExport, error)
}
