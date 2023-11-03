package middleware

// import (
// 	"net/http"
// 	"time"

// 	logging "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
// )

// func LoggerMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		startTime := time.Now()

// 		// Log the start of the request
// 		logging.GetLogger().Info("Started handling request", "requestID", GetReqID(r.Context()), "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)

// 		// Call the next handler in the chain
// 		next.ServeHTTP(w, r)

// 		// Log the end of the request along with the duration
// 		duration := time.Since(startTime)
// 		logging.GetLogger().Info("Completed handling request", "status", w.Header(), "duration", duration.String())
// 	})
// }
