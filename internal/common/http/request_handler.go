package http

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetIDFromRequest(id string, r *http.Request) (uuid.UUID, error) {
	uuidString := mux.Vars(r)[id]
	parsedUUID, err := uuid.Parse(uuidString)

	if err != nil {
		return parsedUUID, errors.New("invalid uuid parameter")
	}
	return parsedUUID, nil
}

func GetloginFromRequest(id string, r *http.Request) string {
	uuidString := mux.Vars(r)[id]

	return uuidString
}

var ErrUnauthorized = &models.UnathorizedError{}

func GetUserFromRequest(r *http.Request) (*models.User, error) {
	user, ok := r.Context().Value(models.ContextKeyUserType{}).(*models.User)
	if !ok {
		return nil, ErrUnauthorized
	}
	if user == nil {
		return nil, ErrUnauthorized
	}

	return user, nil
}
