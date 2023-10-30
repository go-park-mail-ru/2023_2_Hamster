package middleware

import (
	"context"
	"net/http"

	response "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	userRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions"
)

type Middleware struct {
	ur  userRep.Repository
	su  sessions.Usecase
	log logger.CustomLogger
}

func NewMiddleware(su sessions.Usecase, ur userRep.Repository, log logger.CustomLogger) *Middleware {
	return &Middleware{
		su:  su,
		ur:  ur,
		log: log,
	}
}

func (m *Middleware) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			m.log.Errorf("[middleware] no cookie Authentication")
			response.ErrorResponse(w, http.StatusUnauthorized, err, "missing token unauthorized", m.log)
			return
		}
		session := cookie.Value

		m.log.Info("session : " + session)

		if cookie.Value == "" {
			m.log.Errorf("[middleware] missing token")
			response.ErrorResponse(w, http.StatusUnauthorized, err, "missing token unauthorized", m.log) // missing token
			return
		}

		userSession, err := m.su.GetSessionByCookie(context.TODO(), session)
		if err != nil {
			m.log.Errorf("[middleware] validation error: %s", err.Error())
			response.ErrorResponse(w, http.StatusUnauthorized, err, "token validation failed unauthorized", m.log) // token check failed
			return
		}

		user, err := m.ur.GetByID(r.Context(), userSession.UserId)
		if err != nil {
			m.log.Infof("[middleware] get user error: %s", err.Error())
			response.ErrorResponse(w, http.StatusUnauthorized, err, "userAuth check failed", m.log) // UserAuth data check failed
			return
		}

		m.log.Infof("user accepted : %s", user.ID.String())
		reqWithUser := WrapUser(r, user) // Empty user
		next.ServeHTTP(w, reqWithUser)   // token check successed
	})
}

func WrapUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), models.ContextKeyUserType{}, user)
	return r.WithContext(ctx)
}
