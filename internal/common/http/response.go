package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
)

type Response struct {
	Status string      `json:"status"`
	Body   interface{} `json:"body,omitempty"`
}

type Error struct {
	errMsg string `json:"error"`
}

type NilBody struct{}

func NIL() NilBody {
	return NilBody{}
}

func ErrorResponse(w http.ResponseWriter, message string, code int, log logger.CustomLogger) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	errorMsg := Error{
		errMsg: message,
	}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(errorMsg); err != nil {
		log.Errorf("Error failed to marshal error message: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Can't encode error message into json, massage: " + message))
	}
}

func SuccessResponse(w http.ResponseWriter, response any, log logger.CustomLogger) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		log.Errorf("Error failed to marshal response: %s", err.Error())
		ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError, log)
	}
}

func JSON(w http.ResponseWriter, status int, response any) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		w.WriteHeader(status)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// w.Header().Set("Content-Length", )
	w.WriteHeader(status)
}
