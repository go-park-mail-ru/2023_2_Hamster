package app

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/router"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/middleware"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/redis/go-redis/v9"

	authDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/http"
	authRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/repository/postgresql"
	authUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/usecase"

	categoryDelivary "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/delivery/http"
	categoryRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/repository/postgres"
	categoryUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/usecase"

	csrfDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/csrf/delivery/http"
	csrfUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/csrf/usecase"
	transactionDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/transaction/delivery/http"
	transactionRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/transaction/repository/postgresql"
	transactionUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/transaction/usecase"

	userDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/delivery/http"
	userRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/repository/postgresql"
	userUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/usecase"

	sessionRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions/repository/redis"
	sessionUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions/usecase"

	"github.com/gorilla/mux"
)

func Init(db *pgxpool.Pool, redis *redis.Client, log *logger.Logger) *mux.Router {

	authRep := authRep.NewRepository(db, *log)
	sessionRep := sessionRep.NewSessionRepository(redis)

	userRep := userRep.NewRepository(db, *log)
	transactionRep := transactionRep.NewRepository(db, *log)
	categoryRep := categoryRep.NewRepository(db, *log)

	authUsecase := authUsecase.NewUsecase(authRep, userRep, *log)
	sessionUsecase := sessionUsecase.NewSessionUsecase(sessionRep)
	userUsecase := userUsecase.NewUsecase(userRep, *log)
	transactionUsecase := transactionUsecase.NewUsecase(transactionRep, *log)
	categoryUsecase := categoryUsecase.NewUsecase(categoryRep, *log)
	csrfUsecase := csrfUsecase.NewUsecase(*log)

	authMiddlewear := middleware.NewAuthMiddleware(sessionUsecase, userRep, *log)
	logMiddlewear := middleware.NewLoggingMiddleware(*log)
	recoveryMiddlewear := middleware.NewRecoveryMiddleware(*log)
	csrfMiddlewear := middleware.NewCSRFMiddleware(csrfUsecase, *log)

	authHandler := authDelivery.NewHandler(authUsecase, userUsecase, sessionUsecase, *log)
	userHandler := userDelivery.NewHandler(userUsecase, *log)
	transactionHandler := transactionDelivery.NewHandler(transactionUsecase, *log)
	categoryHandler := categoryDelivary.NewHandler(categoryUsecase, *log)
	csrfHandler := csrfDelivery.NewHandler(csrfUsecase, *log)

	return router.InitRouter(
		authHandler,
		userHandler,
		transactionHandler,
		categoryHandler,
		csrfHandler,
		logMiddlewear,
		recoveryMiddlewear,
		authMiddlewear,
		csrfMiddlewear,
	)

}
