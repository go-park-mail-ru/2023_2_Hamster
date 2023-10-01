package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth"
)

type Handler struct {
	au  auth.Usecase
	log logger.CustomLogger
}

func NewHandler(log logger.CustomLogger, au auth.Usecase) *Handler {
	return &Handler{
		au:  au,
		log: log,
	}
}

type signInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if errors.Is(err, io.EOF) {
		h.log.Error(err.Error())
		commonHttp.ErrorResponse(w, "Error request body is empty", http.StatusBadRequest, h.log)
		return
	}
	if err != nil {
		h.log.Errorf("Error decode in request body failed: %v", err.Error())
		commonHttp.ErrorResponse(w, "Error decode in request body failed", http.StatusBadRequest, h.log)
	}

	h.log.Debug("request body successfully decoded", r)

	user := &models.User{}
	if err = json.Unmarshal(body, user); err != nil {
		h.log.Errorf("Error failed to unmarshal request body: %v", err.Error())
		commonHttp.ErrorResponse(w, "failed to unmarshal request body", http.StatusBadRequest, h.log)
		return
	}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) LogOut(w http.ResponseWriter, r *http.Request) {

}
