package middleware

import (
	"context"
	"fmt"
	"net/http"

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
			// handle error
		}
		fmt.Println("-------------------------------------------------")
		fmt.Println(">>>>>>>>>>>", cookie, "<<<<<<<<<<<<<")
		reqToken := &cookie.Raw

		m.log.Info("auth token : " + *reqToken)

		if cookie.Value == "" {
			m.log.Info("[middleware] missing token")
			next.ServeHTTP(w, r) // missing token
			return
		}

		userId, err := m.au.ValidateAccessToken(*reqToken)
		if err != nil {
			m.log.Infof("[middleware] %s", err.Error())
			next.ServeHTTP(w, r) // token check failed
			return
		}

		user, err := m.au.GetUserByAuthData(r.Context(), userId)
		if err != nil {
			m.log.Infof("[middleware] %s", err.Error())
			next.ServeHTTP(w, r) // UserAuth data check failed
			return
		}

		m.log.Infof("user accepted : %d", user.ID)

		ctx := context.WithValue(r.Context(), models.ContextKeyUserType{}, user)
		next.ServeHTTP(w, r.WithContext(ctx)) // token check successed
	})
}
