package router

import (
	//auth "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth/delivery/http"
	"encoding/json"
	"net/http"

	_ "github.com/go-park-mail-ru/2023_2_Hamster/docs"
	user "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/delivery/http"
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
func InitRouter( /*auth *auth.Handler,*/ user *user.Handler /*mid *middleware.Middleware*/) *mux.Router {
	r := mux.NewRouter()

	http.Handle("/", r)
	r.Path("/ping").HandlerFunc(HelloHandler)

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	apiRouter := r.PathPrefix("/api").Subrouter()
	//apiRouter.Use(corsmiddleware.CorsMiddleware)

	// authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	// {
	// 	authRouter.Methods("POST").Path("/signin").HandlerFunc(auth.SignIn)
	// 	authRouter.Methods("POST").Path("/signup").HandlerFunc(auth.SignUp)
	// 	authRouter.Methods("POST").Path("/checkAuth").HandlerFunc(auth.AccessVerification)
	// 	authRouter.Methods("POST").Path("/logout").HandlerFunc(auth.LogOut)
	// }

	userRouter := apiRouter.PathPrefix("/user/{userID}").Subrouter()
	//userRouter.Use(mid.Authentication)
	{
		userRouter.Methods("GET").Path("/").HandlerFunc(user.Get)
		userRouter.Methods("GET").Path("/balance").HandlerFunc(user.GetUserBalance)
		userRouter.Methods("GET").Path("/plannedBudget").HandlerFunc(user.GetPlannedBudget)
		userRouter.Methods("GET").Path("/actualBudget").HandlerFunc(user.GetCurrentBudget)
		userRouter.Methods("GET").Path("/account/all").HandlerFunc(user.GetAccounts)
		userRouter.Methods("GET").Path("/feed").HandlerFunc(user.GetFeed)
	}

	return r
}
