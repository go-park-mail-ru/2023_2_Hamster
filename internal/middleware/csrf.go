package middleware

import (
	"errors"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/csrf"
)

const csrfTokenHttpHeader = "X-CSRF-Token"

type CSRFMiddleware struct {
	tokenServices csrf.Usecase
	logger        logger.Logger
}

func NewCSRFMiddleware(t csrf.Usecase, l logger.Logger) *CSRFMiddleware {
	return &CSRFMiddleware{
		tokenServices: t,
		logger:        l,
	}
}

const (
	missingCSRFToken = "missing CSRF token"
	invalidCSRFToken = "invalid CSRF token"
)

func (m *CSRFMiddleware) CheckCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" {
			user, err := commonHttp.GetUserFromRequest(r)
			if err != nil && errors.Is(err, commonHttp.ErrUnauthorized) {
				commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, commonHttp.ErrUnauthorized.Error(), m.logger)
				return
			}

			clientToken := r.Header.Get(csrfTokenHttpHeader)
			if clientToken == "" {
				commonHttp.ErrorResponse(w, http.StatusBadRequest, errors.New(missingCSRFToken), missingCSRFToken, m.logger)
				return
			}

			userIDFromToken, err := m.tokenServices.CheckCSRFToken(clientToken)
			if err != nil {
				commonHttp.ErrorResponse(w, http.StatusBadRequest, err, invalidCSRFToken, m.logger)
				return
			}
			if user.ID != userIDFromToken {
				commonHttp.ErrorResponse(w, http.StatusBadRequest, err, invalidCSRFToken, m.logger)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
