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

func (u *Usecase) GetBalance(userID uuid.UUID) (float64, error) {
	balance, err := u.userRepo.GetUserBalanceByID(userID)
	if err != nil {
		return 0, fmt.Errorf("(usecase) cant't get balance from repository %w", err)
	}

	return balance, nil
}
