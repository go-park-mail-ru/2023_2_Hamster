package http

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetIDFromRequest(id string, r *http.Request) (uuid.UUID, error) {
	uuidString := mux.Vars(r)[id]
	parsedUUID, err := uuid.Parse(uuidString)

	if err != nil {
		return parsedUUID, errors.New("invalid uuid parametr")
	}
	return parsedUUID, nil
}
