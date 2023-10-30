package usecase

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category"
)

type Usecase struct {
	cu  category.Repository
	log logger.CustomLogger
}

func NewUsecase(cu category.Repository, log logger.CustomLogger)
