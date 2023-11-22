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
