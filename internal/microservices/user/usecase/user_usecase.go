package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user"
	tranfer_models "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/delivery/http/transfer_models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"

	"github.com/google/uuid"
)

type Usecase struct {
	userRepo    user.Repository
	logger      logger.Logger
	accountRepo account.Repository
}

func NewUsecase(
	ur user.Repository,
	log logger.Logger, ar account.Repository) *Usecase {
	return &Usecase{
		userRepo:    ur,
		logger:      log,
		accountRepo: ar,
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
	var errNoSuchAccounts *models.NoSuchAccounts

	dataTranfer.Balance, err = u.GetUserBalance(ctx, userID)
	if err != nil {
		return dataTranfer, err
	}

	dataTranfer.AccountMas, err = u.GetAccounts(ctx, userID)
	if err != nil && !errors.As(err, &errNoSuchAccounts) {
		return dataTranfer, err
	}

	dataTranfer.BudgetPlanned, err = u.GetPlannedBudget(ctx, userID)
	if err != nil {
		return dataTranfer, err
	}

	dataTranfer.BudgetActual, err = u.GetCurrentBudget(ctx, userID)
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

func (u *Usecase) AddUserInAccount(ctx context.Context, accountInput models.AddUserAccount, adminID uuid.UUID) error {
	err := u.accountRepo.SharingCheck(ctx, accountInput.AccountID, adminID)
	if err != nil {
		return fmt.Errorf("[usecase] is not admin user: %w", err)
	}

	sharedUser, err := u.userRepo.GetUserByLogin(ctx, accountInput.Login)
	if err != nil {
		return fmt.Errorf("[usecase] can't get user in login: %w", err)
	}

	err = u.accountRepo.CheckDuplicate(ctx, sharedUser.ID, accountInput.AccountID)
	if err != nil {
		return fmt.Errorf("[usecase] %w", err)
	}

	err = u.accountRepo.AddUserInAccount(ctx, sharedUser.ID, accountInput.AccountID)
	if err != nil {
		return fmt.Errorf("[usecase] can't add user in account %w", err)
	}

	return nil
}

func (u *Usecase) Unsubscribe(ctx context.Context, accountID uuid.UUID, userID uuid.UUID) error {
	err := u.accountRepo.DeleteUserInAccount(ctx, userID, accountID)
	if err != nil {
		return fmt.Errorf("[usecase] can't delete user in account: %w", err)
	}
	return nil
}

func (u *Usecase) DeleteUserInAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, adminID uuid.UUID) error {
	err := u.accountRepo.SharingCheck(ctx, accountID, adminID)
	if err != nil {
		return fmt.Errorf("[usecase] is not admin user: %w", err)
	}

	err = u.accountRepo.DeleteUserInAccount(ctx, userID, accountID)
	if err != nil {
		return fmt.Errorf("[usecase] can't delete account %w", err)
	}

	return nil
}
