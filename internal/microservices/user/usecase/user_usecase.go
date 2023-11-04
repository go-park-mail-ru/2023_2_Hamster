package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user"
	tranfer_models "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/delivery/http/transfer_models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"

	"github.com/google/uuid"
)

type Usecase struct {
	userRepo user.Repository
	logger   logger.Logger
}

func NewUsecase(
	ur user.Repository,
	log logger.Logger) *Usecase {
	return &Usecase{
		userRepo: ur,
		logger:   log,
	}
}

// func (u *Usecase) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error) { // need test
// 	user, err := u.userRepo.GetByID(ctx, userID)
// 	if err != nil {
// 		return user, fmt.Errorf("[usecase] can't get user from repository %w", err)
// 	}

// 	return user, nil
// }

func (u *Usecase) GetUserBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	balance, err := u.userRepo.GetUserBalance(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("[usecase] can't get balance from repository %w", err)
	}

	return balance, nil
}

func (u *Usecase) GetPlannedBudget(ctx context.Context, userID uuid.UUID) (float64, error) {
	balance, err := u.userRepo.GetPlannedBudget(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("[usecase] can't get planned budget from repository %w", err)
	}

	return balance, nil
}

func (u *Usecase) GetCurrentBudget(ctx context.Context, userID uuid.UUID) (float64, error) {
	transactionExpenses, err := u.userRepo.GetCurrentBudget(ctx, userID)

	if err != nil {
		return 0, fmt.Errorf("[usecase] can't get current budget from repository %w", err)
	}

	plannedBudget, err := u.userRepo.GetPlannedBudget(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("[usecase] can't get planned budget from repository %w", err)
	}

	currentBudget := plannedBudget - transactionExpenses
	return currentBudget, nil
}

func (u *Usecase) GetAccounts(ctx context.Context, userID uuid.UUID) ([]models.Accounts, error) { // TO DO MODELS TRANSFER
	account, err := u.userRepo.GetAccounts(ctx, userID)
	if err != nil {
		return account, fmt.Errorf("[usecase] can't get accounts from repository %w", err)
	}

	return account, nil
}

func (u *Usecase) GetFeed(ctx context.Context, userID uuid.UUID) (*tranfer_models.UserFeed, error) { // need test!
	dataTranfer := &tranfer_models.UserFeed{}
	var err error

	dataTranfer.Balance, err = u.GetUserBalance(ctx, userID)
	if err != nil {
		return dataTranfer, err
	}

	dataTranfer.BudgetActual, err = u.GetCurrentBudget(ctx, userID)
	if err != nil {
		return dataTranfer, err
	}

	dataTranfer.BudgetPlanned, err = u.GetPlannedBudget(ctx, userID)
	if err != nil {
		return dataTranfer, err
	}

	dataTranfer.AccountMas, err = u.GetAccounts(ctx, userID)
	if err != nil {
		return dataTranfer, err
	}

	return dataTranfer, nil
}

func (u *Usecase) UpdateUser(ctx context.Context, user *models.User) error { // need test
	if err := u.userRepo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("[usecase] can't update user %w", err)
	}
	return nil
}

func (u *Usecase) UpdatePhoto(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	path := uuid.New()
	err := u.userRepo.UpdatePhoto(ctx, userID, path)
	if err != nil {
		return uuid.Nil, err
	}
	return path, nil
}
