package middleware

import (
	"context"
	"net/http"
	"time"

	response "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	userRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions"
)

type AuthMiddleware struct {
	ur  userRep.Repository
	su  sessions.Usecase
	log logger.CustomLogger
}

func NewAuthMiddleware(su sessions.Usecase, ur userRep.Repository, log logger.CustomLogger) *AuthMiddleware {
	return &AuthMiddleware{
		su:  su,
		ur:  ur,
		log: log,
	}
}

func (m *AuthMiddleware) Authentication(next http.Handler) http.Handler {
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

type LogMiddleware struct {
	log logger.CustomLogger
}

func NewLogMiddleware(log logger.CustomLogger) *LogMiddleware {
	return &LogMiddleware{
		log: log,
	}
}

func (m *LogMiddleware) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		ip, err := response.GetIpFromRequest(r)
		if err != nil {
			m.log.Error(err.Error())
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		params := logger.LogFormatterParams{
			TimeStamp:  time.Now(),
			StatusCode: 200,
			Latency:    time.Since(startTime),
			ClientIP:   ip,
			Method:     r.Method,
			Path:       r.URL.RawPath,
		}

		logMsg := logger.DefaultLogFormatter(params)

		m.log.Info(logMsg)

	})
}

// New will create a new middleware handler from a http.Handler.
func New(h http.Handler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
	}
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "chi/middleware context value " + k.name
}
