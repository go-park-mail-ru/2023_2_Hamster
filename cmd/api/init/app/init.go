package app

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/router"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/redis/go-redis/v9"

	sessionRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions/repository/redis"
	sessionUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions/usecase"
	authDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth/delivery/http"
	authRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth/repository/postgresql"
	authUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth/usecase"
	transactionDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/transaction/delivery/http"
	transactionRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/transaction/repository/postgresql"
	transactionUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/transaction/usecase"
	userDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/delivery/http"
	userRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/repository/postgresql"
	userUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/usecase"
	"github.com/gorilla/mux"
)

func Init(db pgxtype.Querier, redis *redis.Client, log *logger.CustomLogger) *mux.Router {
	authRep := authRep.NewRepository(db, *log)
	sessionRep := sessionRep.NewSessionRepository(redis)
	userRep := userRep.NewRepository(db, *log)
	transactionRep := transactionRep.NewRepository(db, *log)

	authUsecase := authUsecase.NewUsecase(authRep, userRep, *log)
	sessionUsecase := sessionUsecase.NewSessionUsecase(sessionRep, *log)
	userUsecase := userUsecase.NewUsecase(userRep, *log)
	transactionUsecase := transactionUsecase.NewUsecase(transactionRep, *log)
	//middlewear := middleware.NewMiddleware(authUsecase, *log)

	authHandler := authDelivery.NewHandler(authUsecase, userUsecase, sessionUsecase, *log)
	userHandler := userDelivery.NewHandler(userUsecase, *log)
	transactionHandler := transactionDelivery.NewHandler(transactionUsecase, *log)

	return router.InitRouter(authHandler, userHandler, transactionHandler /*, middlewear*/)

}
