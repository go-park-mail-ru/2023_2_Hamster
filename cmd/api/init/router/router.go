package router

import (
	auth "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth/delivery/http"
	user "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/delivery/http"
	"github.com/gorilla/mux"
)

// Initialize router and describes all app's endpoints
func InitRouter(auth *auth.Handler, user *user.Handler) *mux.Router {

	r := mux.NewRouter()

	apiRouter := r.PathPrefix("/api").Subrouter()

	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.Methods("POST").Path("/singin").HandlerFunc(auth.SignIn)
	authRouter.Methods("POST").Path("/sighup").HandlerFunc(auth.SignUp)
	authRouter.Methods("GET").Path("/logout").HandlerFunc(auth.LogOut)

	userRouter := apiRouter.PathPrefix("/user").Subrouter()
	userRouter.Methods("GET").Path("/{userID}/balance").HandlerFunc(user.GetBalance)
	userRouter.Methods("GET").Path("/{userID}/plannedBudget").HandlerFunc(user.GetPlannedBudget)
	userRouter.Methods("GET").Path("/{userID}/actualBudget").HandlerFunc(user.ActualBudget)

	return r
}
