package middleware

//
// import (
// 	"fmt"
// 	"net/http"
// 	"runtime/debug"
//
// 	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
// )
//
// func (m *Middleware) Panic() func(next http.Handler) http.Handler {
//
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			defer func() {
// 				if err := recover(); err != nil {
// 					message := fmt.Errorf("PANIC (recovered): %s\n stacktrace:\n%s", err, string(debug.Stack()))
// 					commonHttp.ErrorResponse(w, http.StatusInternalServerError, message, "server unknown error", m.log)
// 				}
// 			}()
//
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }
