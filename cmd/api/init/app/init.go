package app

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/router"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/jackc/pgx/v5"

	userDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/delivery/http"
	userRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/repository/postgresql"
	userUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/usecase"
	"github.com/gorilla/mux"
)

func Init(db *pgx.Conn, log *logger.CustomLogger) *mux.Router {
	//authRep := authRep.NewRepository(db, *log)
	userRep := userRep.NewRepository(db, *log)

	//authUsecase := authUsecase.NewUsecase(authRep, userRep, *log)
	userUsecase := userUsecase.NewUsecase(userRep, *log)

	//middlewear := middleware.NewMiddleware(authUsecase, *log)

	//authHandler := authDelivery.NewHandler(authUsecase, *log)
	userHandler := userDelivery.NewHandler(userUsecase, *log)

	return router.InitRouter( /*authHandler,*/ userHandler /*, middlewear*/)

}
