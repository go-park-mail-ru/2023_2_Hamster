package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user"
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
		return 0, fmt.Errorf("(usecase) cant't get balance from repository %w", err)
	}

	return balance, nil
}

func (u *Usecase) GetPlannedBudget(userID uuid.UUID) (float64, error) {
	balance, err := u.userRepo.GetPlannedBudget(userID)
	if err != nil {
		return 0, fmt.Errorf("(usecase) cant't get planned budget from repository %w", err)
	}

	return balance, nil
}

func (u *Usecase) GetCurrentBudget(userID uuid.UUID) (float64, error) {
	currentTransaction, err := u.userRepo.GetCurrentBudget(userID)

	if err != nil {
		return 0, fmt.Errorf("(usecase) cant't get current budget from repository %w", err)
	}

	plannedBudget, err := u.userRepo.GetPlannedBudget(userID)
	if err != nil {
		return 0, fmt.Errorf("(usecase) cant't get planned budget from repository %w", err)
	}

	currentBudget := plannedBudget - currentTransaction
	return currentBudget, nil
}
