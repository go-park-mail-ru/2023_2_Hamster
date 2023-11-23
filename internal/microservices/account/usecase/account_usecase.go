package usecase

import (
	"context"
	"fmt"

	logging "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Usecase struct {
	accountRepo account.Repository
	logger      logging.Logger
}

func NewUsecase(
	ar account.Repository,
	log logging.Logger) *Usecase {
	return &Usecase{
		accountRepo: ar,
		logger:      log,
	}
}

func (a *Usecase) CreateAccount(ctx context.Context, userID uuid.UUID, account *models.Accounts) (uuid.UUID, error) {
	accountID, err := a.accountRepo.CreateAccount(ctx, userID, account)

	if err != nil {
		return accountID, fmt.Errorf("[usecase] can't create account into repository: %w", err)
	}
	return accountID, nil
}

func (a *Usecase) UpdateAccount(ctx context.Context, userID uuid.UUID, account *models.Accounts) error {
	err := a.accountRepo.CheckForbidden(ctx, account.ID, userID)
	if err != nil {
		return fmt.Errorf("[usecase] can't be update by user: %w", err)
	}

	err = a.accountRepo.UpdateAccount(ctx, userID, account)
	if err != nil {
		return fmt.Errorf("[usecase] can't update account into repository: %w", err)
	}
	return nil
}

func (a *Usecase) DeleteAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error {
	err := a.accountRepo.CheckForbidden(ctx, accountID, userID)
	if err != nil {
		return fmt.Errorf("[usecase] can't be delete by user: %w", err)
	}

	err = a.accountRepo.DeleteAccount(ctx, userID, accountID)
	if err != nil {
		return fmt.Errorf("[usecase] can't create account into repository: %w", err)
	}
	return nil
}
