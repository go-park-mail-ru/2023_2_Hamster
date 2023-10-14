package usecase

import (
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user"
	tranfer_models "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/delivery/http/transfer_models"
	"github.com/hashicorp/go-multierror"

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

func (u *Usecase) GetFeed(userID uuid.UUID) (tranfer_models.UserFeed, *multierror.Error) { // need test!
	var dataTranfer tranfer_models.UserFeed
	var err error
	var multiErr *multierror.Error

	errMsg := "errors: "
	dataTranfer.Balance, err = u.GetUserBalance(userID)
	if err != nil {
		multiErr = multierror.Append(multiErr, errors.New("balance: "+err.Error()))

		errMsg += "(balance) "
	}

	dataTranfer.BudgetActual, err = u.GetCurrentBudget(userID)
	if err != nil {
		multiErr = multierror.Append(multiErr, errors.New("current: "+err.Error()))
		errMsg += "(budgetActual) "
	}

	dataTranfer.BudgetPlanned, err = u.GetPlannedBudget(userID)
	if err != nil {
		multiErr = multierror.Append(multiErr, errors.New("planned: "+err.Error()))
		errMsg += "(budgetPlanned)"
	}

	dataTranfer.Account.Account, err = u.GetAccounts(userID)
	if err != nil {
		multiErr = multierror.Append(multiErr, errors.New("accounts: "+err.Error()))
		errMsg += "(account)"
	}

	dataTranfer.ErrMes = errMsg
	return dataTranfer, multiErr
}
