package http

import (
	"net/http"
	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user"
)

type Handler struct {
	userService user.Usecase
	logger      logger.CustomLogger
}

const (
	userIdUrlParam  = "userID"
)

func NewHandler(uu user.Usecase, l logger.CustomLogger) *Handler {
	return &Handler{
		userService: uu,
		logger:      l,
	}
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	commonHttp.GetIDFromRequest(r ,userIdUrlParam)
	if err != nil {
		h.logger.infof("invalid id: %v:", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

}

func (h* Handler) GetPlannedBudget(w http.ResponseWriter, r *http.Request) {
	commonHttp.
}

func (h* Handler) ActualBudget(w http.ResponseWriter, r *http.Request) [
	commonHttp.
]