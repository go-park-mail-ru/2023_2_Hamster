package app

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/router"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	authDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth/delivery/http"
	authRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth/repository/postgresql"
	authUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth/usecase"
	userDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/delivery/http"
	userRer "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/repository/postgresql"
	userUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/usecase"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func Init(db *sqlx.DB, log logger.CustomLogger) *mux.Router {
	authRep := authRep.NewRepository(db, log)
	userRep := userRer.NewRepository(db, log)

	authUsecase := authUsecase.NewUsecase(authRep, userRep, log)
	userUsecase := userUsecase.NewUsecase(userRep, log)

	authHandler := authDelivery.NewHandler(authUsecase, log)
	userHandler := userDelivery.NewHandler(userUsecase, log)

	return router.InitRouter(
		authHandler,
	)

}
