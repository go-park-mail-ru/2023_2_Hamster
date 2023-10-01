package app

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
)

func Init(db *sqlx.DB, logger logger.CustomLogger) *mux.Router {
	authRepo := authRepository.userRepository(db)

	authUse
}
