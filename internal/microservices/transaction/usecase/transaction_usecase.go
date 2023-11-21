package usecase

import (
	"context"
	"fmt"

	logging "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/transaction"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Usecase struct {
	transactionRepo transaction.Repository
	logger          logging.Logger
}

func NewUsecase(
	tr transaction.Repository,
	log logging.Logger) *Usecase {
	return &Usecase{
		transactionRepo: tr,
		logger:          log,
	}
}

func (t *Usecase) GetFeed(ctx context.Context, userID uuid.UUID, query *models.QueryListOptions) ([]models.Transaction, error) {
	transaction, err := t.transactionRepo.GetFeed(ctx, userID, query)
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

func (t *Usecase) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	userIDCheck, err := t.transactionRepo.CheckForbidden(ctx, transaction.ID)
	if err != nil {
		return fmt.Errorf("[usecase] can't find transaction in repository %w", err)
	}

	if userIDCheck != transaction.UserID {
		return fmt.Errorf("[usecase] can't be update by user: %w", &models.ForbiddenUserError{})
	}

	if err := t.transactionRepo.UpdateTransaction(ctx, transaction); err != nil {
		return fmt.Errorf("[usecase] can't update transaction %w", err)
	}
	return nil
}

func (t *Usecase) DeleteTransaction(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID) error {
	userIDCheck, err := t.transactionRepo.CheckForbidden(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("[usecase] can't find transaction in repository %w", err)
	}

	if userIDCheck != userID {
		return fmt.Errorf("[usecase] can't be deleted by user: %w", &models.ForbiddenUserError{})
	}

	err = t.transactionRepo.DeleteTransaction(ctx, transactionID, userID)
	if err != nil {
		return fmt.Errorf("[usecase] can`t be deleted from repository")
	}

	return nil
}
