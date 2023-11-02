package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category"
	"github.com/google/uuid"
)

type CategoryCreate struct {
	ID   uuid.UUID
	Name string
}

type Handler struct {
	cu  category.Usecase
	log logger.Logger
}

func NewHandler(cu category.Usecase, log logger.Logger) *Handler {
	return &Handler{
		cu:  cu,
		log: log,
	}
}

func (h *Handler) CreateCategry(w http.ResponseWriter, r http.Request) {

}
