package http

import (
	"errors"
	"net/http"
	"strconv"

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

func GetloginFromRequest(login string, r *http.Request) string {
	userLogin := mux.Vars(r)[login]

	return userLogin
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

func GetQueryParam(r *http.Request) (int, int, error) {
	values := r.URL.Query()
	var page, perPage int
	var err error

	if pageStr := values.Get("page"); pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if page < 0 {
			return 0, 0, errors.New("error page < 0")
		}
		if err != nil {
			return 0, 0, err
		}
	}

	if perPageStr := values.Get("page_size"); perPageStr != "" {
		perPage, err = strconv.Atoi(perPageStr)
		if perPage < 0 {
			return 0, 0, errors.New("error page_size < 0")
		}
		if err != nil {
			return 0, 0, err
		}
	}

	return page, perPage, nil
}
