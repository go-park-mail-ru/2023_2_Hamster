package usecase

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user"
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
