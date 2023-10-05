package corsmiddlware

import (
	"net/http"
	"os"
	"strconv"
)

const (
	Origin      = "Access-Control-Allow-Origin"
	Methods     = "Access-Control-Allow-Methods"
	Headers     = "Access-Control-Allow-Headers"
	MaxAge      = "Access-Control-Max-Age"
	Credentials = "Access-Control-Allow-Credentials"

	AllowedMethods = "GET,POST,PUT,DELETE"
	AllowedHeaders = "Content-Type, Cookie"
	MaxAgeValue    = 24 * 60 * 60
)

func corsMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(Origin, os.Getenv("FRONTEND_IP")+os.Getenv("FRONTEND_PORT"))
		w.Header().Set(Methods, AllowedMethods)
		w.Header().Set(Headers, AllowedHeaders)
		w.Header().Set(Credentials, "true")
		w.Header().Set(MaxAge, strconv.Itoa(MaxAgeValue))

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}
