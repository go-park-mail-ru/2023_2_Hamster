package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/transaction"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
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

func (t *Usecase) GetFeed(ctx context.Context, userID uuid.UUID, page int, pageSize int) ([]models.Transaction, bool, error) {
	transaction, isAll, err := t.transactionRepo.GetFeed(ctx, userID, page, pageSize)
	if err != nil {

		return transaction, isAll, fmt.Errorf("[usecase] can't get transactions from repository %w", err)
	}
	return transaction, isAll, nil
}

func (t *Usecase) CreateTransaction(ctx context.Context, transaction *models.Transaction) (uuid.UUID, error) {
	transactionID, err := t.transactionRepo.CreateTransaction(ctx, transaction)

	if err != nil {
		return transactionID, fmt.Errorf("[usecase] can't create transaction into repository: %w", err)
	}
	return transactionID, nil
}

func (t *Usecase) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	if err := t.transactionRepo.UpdateTransaction(ctx, transaction); err != nil {
		return fmt.Errorf("[usecase] can't update transaction %w", err)
	}
	return nil
}

func (t *Usecase) DeleteTransaction(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID) error {
	userIDCheck, err := t.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("[usecase] can't find artist in repository")
	}

	if userIDCheck == userID {
		return fmt.Errorf("[usecase] can't be deleted by user: %w", &models.ForbiddenUserError{})
	}

	err = t.DeleteTransaction(ctx, transactionID, userID)
	if err != nil {
		return fmt.Errorf("[usecase] can`t be deleted from repository")
	}

	return nil
}
