package router

import (
	//auth "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth/delivery/http"

	"net/http"
	"time"

	_ "github.com/go-park-mail-ru/2023_2_Hamster/docs"

	auth "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/http"
	category "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/delivery/http"
	csrf "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/csrf/delivery/http"
	transaction "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/transaction/delivery/http"
	user "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/delivery/http"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/middleware"
	"github.com/gorilla/mux"

	httpSwagger "github.com/swaggo/http-swagger"
)

// Initialize router and describes all app's endpoints
func InitRouter(auth *auth.Handler,
	user *user.Handler,
	transaction *transaction.Handler,
	category *category.Handler,
	csrf *csrf.Handler,
	logMid *middleware.LoggingMiddleware,
	recoveryMid *middleware.RecoveryMiddleware,
	authMid *middleware.AuthMiddleware,
	csrfMid *middleware.CSRFMiddleware) *mux.Router {

	r := mux.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(logMid.LoggingMiddleware)
	r.Use(recoveryMid.Recoverer)
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware.Heartbeat("ping"))

	http.Handle("/", r)

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	apiRouter := r.PathPrefix("/api").Subrouter()

	csrfRouter := apiRouter.PathPrefix("/csrf").Subrouter()
	csrfRouter.Use(authMid.Authentication)
	csrfRouter.Methods("GET").Path("/").HandlerFunc(csrf.GetCSRF)

	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	{
		authRouter.Methods("POST").Path("/signin").HandlerFunc(auth.Login)
		authRouter.Methods("POST").Path("/signup").HandlerFunc(auth.SignUp)
		authRouter.Methods("POST").Path("/checkAuth").HandlerFunc(auth.HealthCheck)
		authRouter.Methods("POST").Path("/loginCheck").HandlerFunc(auth.CheckLoginUnique)
		authRouter.Methods("POST").Path("/logout").HandlerFunc(auth.LogOut)
	}

	userRouter := apiRouter.PathPrefix("/user").Subrouter()
	userRouter.Use(authMid.Authentication)
	userRouter.Use(csrfMid.CheckCSRF)
	{
		userRouter.Methods("PUT").Path("/updatePhoto").HandlerFunc(user.UpdatePhoto)
		userRouter.Methods("PUT").Path("/update").HandlerFunc(user.Update)

		// userRouter.Methods("GET").Path("/balance").HandlerFunc(user.GetUserBalance)
		// userRouter.Methods("GET").Path("/plannedBudget").HandlerFunc(user.GetPlannedBudget)
		// userRouter.Methods("GET").Path("/actualBudget").HandlerFunc(user.GetCurrentBudget)
		userRouter.Methods("GET").Path("/account/all").HandlerFunc(user.GetAccounts)
		userRouter.Methods("GET").Path("/feed").HandlerFunc(user.GetFeed)
		userRouter.Methods("GET").Path("/").HandlerFunc(user.Get)

	}

	transactionRouter := apiRouter.PathPrefix("/transaction").Subrouter()
	transactionRouter.Use(authMid.Authentication)
	transactionRouter.Use(csrfMid.CheckCSRF)
	{
		transactionRouter.Methods("GET").Path("/feed").HandlerFunc(transaction.GetFeed)
		// 	transactionRouter.Methods("GET").Path("/{transaction_id}/").HandlerFunc(transaction.Get)
		transactionRouter.Methods("PUT").Path("/update").HandlerFunc(transaction.Update)
		transactionRouter.Methods("POST").Path("/create").HandlerFunc(transaction.Create)
		transactionRouter.Methods("DELETE").Path("/{transaction_id}/delete").HandlerFunc(transaction.Delete)
	}

	categoryRouter := apiRouter.PathPrefix("/tag").Subrouter()
	categoryRouter.Use(authMid.Authentication)
	categoryRouter.Use(csrfMid.CheckCSRF)
	{
		categoryRouter.Methods("POST").Path("/create").HandlerFunc(category.CreateTag)
		categoryRouter.Methods("GET").Path("/all").HandlerFunc(category.GetTags)
		categoryRouter.Methods("PUT").Path("/{tagID}/update").HandlerFunc(category.UpdateTag)
		categoryRouter.Methods("DELETE").Path("/delete").HandlerFunc(category.DeleteTag)
	}
	return r
}
