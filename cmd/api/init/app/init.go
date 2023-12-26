package app

import (
	"os"

	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/router"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/middleware"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/redis/go-redis/v9"

	generatedAccount "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account/delivery/grpc/generated"
	generatedAuth "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/grpc/generated"
	authDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/http"
	generatedCategory "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/delivery/grpc/generated"

	categoryDelivary "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/delivery/http"

	csrfDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/csrf/delivery/http"
	csrfUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/csrf/usecase"

	transactionDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/transaction/delivery/http"
	transactionRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/transaction/repository/postgresql"
	transactionUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/transaction/usecase"

	userDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/delivery/http"
	userRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/repository/postgresql"
	userUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/usecase"

	accountDelivery "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account/delivery/http"
	accountRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account/repository/postgresql"
	sessionRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions/repository/redis"
	sessionUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions/usecase"

	iohub "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/IOhub/delivery/http"

	"github.com/gorilla/mux"
)

func Init(db *pgxpool.Pool, redis *redis.Client, log *logger.Logger) *mux.Router {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erorr loading .env file %v\n", err)
	}

	authAddr := os.Getenv("AUTH_ADDR")
	accountAddr := os.Getenv("ACCOUNT_ADDR")
	categoryAddr := os.Getenv("CATEGORY_ADDR")

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.FailOnNonTempDialError(true),
	}

	authConn, err := grpc.Dial(authAddr, opts...)
	if err != nil {
		log.Fatalf("Connection refused auth: %v\n", err)
	}

	authClient := generatedAuth.NewAuthServiceClient(authConn)

	accountConn, err := grpc.Dial(accountAddr, opts...)

	if err != nil {
		log.Fatalf("Connection refused account %v\n", err)
	}

	categoryConn, err := grpc.Dial(categoryAddr, opts...)

	if err != nil {
		log.Fatalf("Connection refused category %v\n", err)
	}

	categortClient := generatedCategory.NewCategoryServiceClient(categoryConn)

	accountClient := generatedAccount.NewAccountServiceClient(accountConn)
	// authRep := authRep.NewRepository(db, *log)
	sessionRep := sessionRep.NewSessionRepository(redis)

	userRep := userRep.NewRepository(db, *log)
	transactionRep := transactionRep.NewRepository(db, *log)
	//categoryRep := categoryRep.NewRepository(db, *log)
	accountRep := accountRep.NewRepository(db, *log)

	// authUsecase := authUsecase.NewUsecase(authRep, *log)
	sessionUsecase := sessionUsecase.NewSessionUsecase(sessionRep)
	userUsecase := userUsecase.NewUsecase(userRep, *log, accountRep)
	transactionUsecase := transactionUsecase.NewUsecase(transactionRep, *log)
	//categoryUsecase := categoryUsecase.NewUsecase(categoryRep, *log)
	csrfUsecase := csrfUsecase.NewUsecase(*log)
	// accountUsecase := accountUsecase.NewUsecase(accountRep, *log)

	authHandler := authDelivery.NewHandler(sessionUsecase, authClient, *log)

	authMiddlewear := middleware.NewAuthMiddleware(sessionUsecase, userRep, *log)
	logMiddlewear := middleware.NewLoggingMiddleware(*log)
	recoveryMiddlewear := middleware.NewRecoveryMiddleware(*log)
	csrfMiddlewear := middleware.NewCSRFMiddleware(csrfUsecase, *log)

	userHandler := userDelivery.NewHandler(userUsecase, *log)
	transactionHandler := transactionDelivery.NewHandler(transactionUsecase, userUsecase, accountClient, *log)
	categoryHandler := categoryDelivary.NewHandler(categortClient, *log)
	csrfHandler := csrfDelivery.NewHandler(csrfUsecase, *log)
	accountHandler := accountDelivery.NewHandler(accountClient, *log)
	iohubHandler := iohub.NewHandler(transactionUsecase, userUsecase, categortClient, accountClient, *log)

	return router.InitRouter(
		authHandler,
		userHandler,
		transactionHandler,
		categoryHandler,
		accountHandler,
		iohubHandler,
		csrfHandler,
		logMiddlewear,
		recoveryMiddlewear,
		authMiddlewear,
		csrfMiddlewear,
	)

}
