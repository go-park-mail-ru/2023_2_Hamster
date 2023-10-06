package corsmiddleware

import (
	"net/http"
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

func CorsMiddleware(next http.Handler) http.Handler {
	//origin := "http://" + os.Getenv("FRONTEND_IP") + ":" + os.Getenv("FRONTEND_PORT") + ", " + "http://127.0.0.1:8000"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(Origin /*, origin */, "*")
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
