package router

import (
	//auth "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth/delivery/http"
	"encoding/json"
	"net/http"

	_ "github.com/go-park-mail-ru/2023_2_Hamster/docs"
	auth "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/http"
	csrf "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/csrf/delivery/http"
	transaction "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/transaction/delivery/http"
	user "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/delivery/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/middleware"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status string
		Body   interface{}
	}{
		Status: "200",
		Body:   "Pong",
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Initialize router and describes all app's endpoints
func InitRouter(auth *auth.Handler,
	user *user.Handler,
	transaction *transaction.Handler,
	csrf *csrf.Handler,
	authMid *middleware.AuthMiddleware,
	csrfMid *middleware.CSRFMiddleware,
	/*logMid *middleware.LogMiddleware*/) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.RequestID)
	// r.Use(logMid.LoggerMiddleware)
	//r.Use(middleware.Logger)
	// r.Use(mid.Panic())

	http.Handle("/", r)

	r.Path("/ping").HandlerFunc(HelloHandler)
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
		authRouter.Methods("GET").Path("/check-unique-login/{login}").HandlerFunc(auth.CheckLoginUnique)
		authRouter.Methods("POST").Path("/logout").HandlerFunc(auth.LogOut)
	}

	userRouter := apiRouter.PathPrefix("/user").Subrouter()
	userRouter.Use(authMid.Authentication)
	userRouter.Use(csrfMid.CheckCSRF)
	{
		userRouter.Methods("PUT").Path("/updatePhoto").HandlerFunc(user.UpdatePhoto)
		userRouter.Path("/update").Methods("PUT").HandlerFunc(user.Update)

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

	// categoryRouter := apiRouter.PathPrefix("/category").Subrouter()
	// {
	// 	categoryRouter.Methods("GET").Path("/all").HandlerFunc(category.GetFeed)
	// 	categoryRouter.Methods("PUT").Path("/{categoryID}/update").HandlerFunc(category.Update)
	// 	categoryRouter.Methods("POST").Path("/create").HandlerFunc(category.Create)
	// 	categoryRouter.Methods("DELETE").Path("/delete").HandlerFunc(category.Delete)
	// }
	return r
}
