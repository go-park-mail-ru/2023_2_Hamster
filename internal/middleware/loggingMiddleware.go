package middleware

import (
	"net/http"
	"time"

	contextutils "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/context_utils"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/sirupsen/logrus"
)

type ResponseWriterWrap struct {
	http.ResponseWriter
	Status int
	Length int
}

func (r *ResponseWriterWrap) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *ResponseWriterWrap) Write(bytes []byte) (int, error) {
	r.Length = len(bytes)

	return r.ResponseWriter.Write(bytes)
}

type LoggingMiddleware struct {
	log logger.Logger
}

func NewLoggingMiddleware(log logger.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		log: log,
	}
}

func (m *LoggingMiddleware) LoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wr := &ResponseWriterWrap{
			ResponseWriter: w,
			Status:         200,
		}

		next.ServeHTTP(wr, r)

		status := wr.Status
		length := wr.Length

		midlog := m.log.WithFields(logrus.Fields{
			"time":       time.Now(),
			"duration":   time.Since(startTime),
			"Request-ID": contextutils.GetReqID(r.Context()),
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     status,
			"remote-IP":  r.RemoteAddr,
			"byteLen":    length,
			"user-agent": r.UserAgent(),
		})
		switch {
		case status >= http.StatusInternalServerError:
			midlog.Error("Server Error")
		case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
			midlog.Warn("Client Error")
		case status >= http.StatusMultipleChoices && status < http.StatusBadRequest:
			midlog.Info("Redirect")
		case status >= http.StatusOK && status < http.StatusMultipleChoices:
			midlog.Info("Success")
		default:
			midlog.Info("Informational")
		}
	}
	return http.HandlerFunc(fn)
}
