package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user"
	tranfer_models "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/delivery/http/transfer_models"

	"github.com/google/uuid"
)

type Usecase struct {
	userRepo user.Repository
	logger   logger.CustomLogger
}

func NewUsecase(
	ur user.Repository,
	log logger.CustomLogger) *Usecase {
	return &Usecase{
		userRepo: ur,
		logger:   log,
	}
}

func (u *Usecase) GetUser(userID uuid.UUID) (*models.User, error) { // need test
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return user, fmt.Errorf("[usecase] can't get user from repository %w", err)
	}

	return user, nil
}

func (u *Usecase) GetUserBalance(userID uuid.UUID) (float64, error) {
	balance, err := u.userRepo.GetUserBalance(userID)
	if err != nil {
		return 0, fmt.Errorf("[usecase] can't get balance from repository %w", err)
	}

	return balance, nil
}

func (u *Usecase) GetPlannedBudget(userID uuid.UUID) (float64, error) {
	balance, err := u.userRepo.GetPlannedBudget(userID)
	if err != nil {
		return 0, fmt.Errorf("[usecase] can't get planned budget from repository %w", err)
	}

	return balance, nil
}

func (u *Usecase) GetCurrentBudget(userID uuid.UUID) (float64, error) { // need test
	transactionExpenses, err := u.userRepo.GetCurrentBudget(userID)

	if err != nil {
		return 0, fmt.Errorf("[usecase] can't get current budget from repository %w", err)
	}

	plannedBudget, err := u.userRepo.GetPlannedBudget(userID)
	if err != nil {
		return 0, fmt.Errorf("[usecase] can't get planned budget from repository %w", err)
	}

	currentBudget := plannedBudget - transactionExpenses
	return currentBudget, nil
}

func (u *Usecase) GetAccounts(userID uuid.UUID) ([]models.Accounts, error) {
	account, err := u.userRepo.GetAccounts(userID)
	if err != nil {

		return account, fmt.Errorf("[usecase] can't get accounts from repository %w", err)
	}

	return account, nil
}

func (u *Usecase) GetFeed(userID uuid.UUID) (tranfer_models.UserFeed, error) { // need test!
	var dataTranfer tranfer_models.UserFeed
	var err error

	dataTranfer.Balance, err = u.GetUserBalance(userID)
	if err != nil {

		return dataTranfer, err
	}

	dataTranfer.BudgetActual, err = u.GetCurrentBudget(userID)
	if err != nil {
		return dataTranfer, err
	}

	dataTranfer.BudgetPlanned, err = u.GetPlannedBudget(userID)
	if err != nil {
		return dataTranfer, err
	}

	dataTranfer.AccountMas, err = u.GetAccounts(userID)
	if err != nil {
		return dataTranfer, err
	}

	return dataTranfer, nil
}
