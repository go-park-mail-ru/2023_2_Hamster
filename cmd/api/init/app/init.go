package app

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/router"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/redis/go-redis/v9"

	generatedAccount "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account/delivery/grpc/generated"
	generatedAuth "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/grpc/generated"
	authDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/http"

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

	accountDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account/delivery/http"

	sessionRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions/repository/redis"
	sessionUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions/usecase"

	"github.com/gorilla/mux"
)

func Init(db *pgxpool.Pool, redis *redis.Client, log *logger.Logger) *mux.Router {
	authConn, err := grpc.Dial(
		"0.0.0.0:8010",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		err = fmt.Errorf("error happened in grpc.Dial auth: %w", err)
		return nil
	}
	defer authConn.Close()

	authClient := generatedAuth.NewAuthServiceClient(authConn)

	// authRep := authRep.NewRepository(db, *log)
	accountConn, err := grpc.Dial(
		"127.0.0.1:8030",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatal("can't connect to account grpc")
	}

	accountClient := generatedAccount.NewAccountServiceClient(accountConn)
	// authRep := authRep.NewRepository(db, *log)
	sessionRep := sessionRep.NewSessionRepository(redis)

	userRep := userRep.NewRepository(db, *log)
	transactionRep := transactionRep.NewRepository(db, *log)
	categoryRep := categoryRep.NewRepository(db, *log)
	// accountRep := accountRep.NewRepository(db, *log)

	// authUsecase := authUsecase.NewUsecase(authRep, *log)
	sessionUsecase := sessionUsecase.NewSessionUsecase(sessionRep)
	userUsecase := userUsecase.NewUsecase(userRep, *log)
	transactionUsecase := transactionUsecase.NewUsecase(transactionRep, *log)
	categoryUsecase := categoryUsecase.NewUsecase(categoryRep, *log)
	csrfUsecase := csrfUsecase.NewUsecase(*log)
	// accountUsecase := accountUsecase.NewUsecase(accountRep, *log)

	authHandler := authDelivery.NewHandler(sessionUsecase, authClient, *log)

	authMiddlewear := middleware.NewAuthMiddleware(sessionUsecase, userRep, *log)
	logMiddlewear := middleware.NewLoggingMiddleware(*log)
	recoveryMiddlewear := middleware.NewRecoveryMiddleware(*log)
	csrfMiddlewear := middleware.NewCSRFMiddleware(csrfUsecase, *log)

	userHandler := userDelivery.NewHandler(userUsecase, *log)
	transactionHandler := transactionDelivery.NewHandler(transactionUsecase, *log)
	categoryHandler := categoryDelivary.NewHandler(categoryUsecase, *log)
	csrfHandler := csrfDelivery.NewHandler(csrfUsecase, *log)
	accountHandler := accountDelivery.NewHandler(accountClient, *log)

	return router.InitRouter(
		authHandler,
		userHandler,
		transactionHandler,
		categoryHandler,
		csrfHandler,
		accountHandler,
		logMiddlewear,
		recoveryMiddlewear,
		authMiddlewear,
		csrfMiddlewear,
	)

}
