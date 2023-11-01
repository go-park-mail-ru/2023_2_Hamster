package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
)

const (
	InvalidURLParameter = "invalid url parameter"
	InvalidBodyRequest  = "invalid input body"
	ForbiddenUser       = "user has no rights"
)

const minErrorToLogCode = 500

type Response[T any] struct {
	Status int `json:"status"`
	Body   T   `json:"body"`
}

type ResponseError struct {
	Status int    `json:"status"`
	ErrMes string `json:"message"`
}

type NilBody struct{}

func ErrorResponse(w http.ResponseWriter, code int, err error, message string, log logger.CustomLogger) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	errorMsg := ResponseError{
		Status: code,
		ErrMes: message,
	}

	if code < minErrorToLogCode {
		log.Infof("invalid request: %v:", err)
	} else {
		log.Error(err.Error())
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(errorMsg); err != nil {
		log.Errorf("Error failed to marshal error message: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		if _, writeErr := w.Write([]byte("Can't encode error message into json, message: " + message)); writeErr != nil {
			log.Errorf("Error writing response: %s", writeErr.Error())
		}
	}
}

func SuccessResponse[T any](w http.ResponseWriter, status int, response T) {
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
