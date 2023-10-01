package router

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth"
	"github.com/gorilla/mux"
)

// Initialize router and describes all app's endpoints
func InitRouter() {
	authRepo := 

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/signup", )
	}
}
