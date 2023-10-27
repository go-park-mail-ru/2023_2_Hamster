package app

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/router"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/jackc/pgtype/pgxtype"

	transactionDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/transaction/delivery/http"
	transactionRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/transaction/repository/postgresql"
	transactionUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/transaction/usecase"
	userDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/delivery/http"
	userRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/repository/postgresql"
	userUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/usecase"
	"github.com/gorilla/mux"
)

func Init(db pgxtype.Querier, log *logger.CustomLogger) *mux.Router {
	//authRep := authRep.NewRepository(db, *log)
	userRep := userRep.NewRepository(db, *log)
	transactionRep := transactionRep.NewRepository(db, *log)

	//authUsecase := authUsecase.NewUsecase(authRep, userRep, *log)
	userUsecase := userUsecase.NewUsecase(userRep, *log)
	transactionUsecase := transactionUsecase.NewUsecase(transactionRep, *log)
	//middlewear := middleware.NewMiddleware(authUsecase, *log)

	//authHandler := authDelivery.NewHandler(authUsecase, *log)
	userHandler := userDelivery.NewHandler(userUsecase, *log)
	transactionHandler := transactionDelivery.NewHandler(transactionUsecase, *log)

	return router.InitRouter( /*authHandler,*/ userHandler, transactionHandler /*, middlewear*/)

}
