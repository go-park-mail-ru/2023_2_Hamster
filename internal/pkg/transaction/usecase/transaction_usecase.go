package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/transaction"
	"github.com/google/uuid"
)

type Usecase struct {
	transactionRepo transaction.Repository
	logger          logger.CustomLogger
}

func NewUsecase(
	tr transaction.Repository,
	log logger.CustomLogger) *Usecase {
	return &Usecase{
		transactionRepo: tr,
		logger:          log,
	}
}

func (t *Usecase) GetFeed(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	transaction, err := t.transactionRepo.GetFeed(ctx, userID)
	if err != nil {

		return transaction, fmt.Errorf("[usecase] can't get transactions from repository %w", err)
	}
	return transaction, nil
}

func (t *Usecase) CreateTransaction(ctx context.Context, transaction *models.Transaction) (uuid.UUID, error) {
	transactionID, err := t.transactionRepo.CreateTransaction(ctx, transaction)

	if err != nil {
		return transactionID, fmt.Errorf("[usecase] can't create transaction into repository: %w", err)
	}
	return transactionID, nil
}
