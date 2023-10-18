package middleware

import (
	"context"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth"
)

type Middleware struct {
	au  auth.Usecase
	log logger.CustomLogger
}

func NewMiddleware(au auth.Usecase, log logger.CustomLogger) *Middleware {
	return &Middleware{
		au:  au,
		log: log,
	}
}

func (m *Middleware) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Authentication")
		if err != nil {
			m.log.Errorf("[middleware] no cookie Authentication")
			commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, "missing token unauthorized", m.log)
			// next.ServeHTTP(w, r)
			return
		}
		reqToken := cookie.Value

		m.log.Info("auth token : " + reqToken)

		if cookie.Value == "" {
			m.log.Errorf("[middleware] missing token")
			commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, "missing token unauthorized", m.log) // missing token
			return
		}

		userId, _, err := m.au.ValidateAccessToken(reqToken)
		if err != nil {
			m.log.Errorf("[middleware] validation error: %s", err.Error())
			commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, "token validation failed unauthorized", m.log) // token check failed
			return
		}

		//fmt.Println(userId)
		user, err := m.au.GetUserByAuthData(r.Context(), userId)
		if err != nil {
			m.log.Infof("[middleware] get user error: %s", err.Error())
			commonHttp.ErrorResponse(w, http.StatusUnauthorized, err, "userAuth check failed", m.log) // UserAuth data check failed
			return
		}

		m.log.Infof("user accepted : %d", user.ID)

		ctx := context.WithValue(r.Context(), models.ContextKeyUserType{}, user)
		next.ServeHTTP(w, r.WithContext(ctx)) // token check successed
	})
}
