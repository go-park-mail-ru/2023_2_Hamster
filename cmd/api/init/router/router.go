package router

import "github.com/gorilla/mux"

// Initialize router and describes all app's endpoints
func InitRouter() {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/signup")
	}
}
