package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
)

type RecoveryMiddleware struct {
	log logger.Logger
}

func NewRecoveryMiddleware(log logger.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		log: log,
	}
}

func (m *RecoveryMiddleware) Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					panic(rvr)
				}

				m.log.Panic(rvr, debug.Stack())

				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
