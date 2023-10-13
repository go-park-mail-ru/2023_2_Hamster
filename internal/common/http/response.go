package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
)

type Error struct {
	ErrMes string `json:"message"`
}

type NilBody struct{}

type Response[T any] struct {
	Status int `json:"status"`
	Body   T   `json:"body"`
}

func NIL() NilBody {
	return NilBody{}
}

func ErrorResponse(w http.ResponseWriter, code int, message string, log logger.CustomLogger) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	errorMsg := Error{
		ErrMes: message,
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(errorMsg); err != nil {
		log.Errorf("Error failed to marshal error message: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Can't encode error message into json, massage: " + message))
	}
}

func JSON[T any](w http.ResponseWriter, status int, response T) {
	date := Response[T]{Status: status, Body: response}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(date); err != nil {
		w.WriteHeader(status)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// w.Header().Set("Content-Length", )
	w.WriteHeader(status)
}
