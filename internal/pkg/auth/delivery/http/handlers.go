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

type signUpResponse struct {
	ID uint32 `json:"id"`
}

type signInput struct {
	Username string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required"`
}

type loginResponse struct {
	JWT string `json:"jwt"`
}

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

	id, err := h.au.SignUpUser(r.Context(), user)
	var errUserAlreadyExist *models.UserAlreadyExistsError
	if errors.As(err, errUserAlreadyExist) {
		h.log.Error(err.Error())
		commonHttp.ErrorResponse(w, "user already exist", http.StatusBadRequest, h.log)
		return
	} else if err != nil {
		h.log.Error(err.Error())
		commonHttp.ErrorResponse(w, "server error", http.StatusBadRequest, h.log)
		return
	}

	h.log.Infof("User created with id: %d", id)

	suResp := signUpResponse{ID: id}

	commonHttp.SuccessResponse(w, suResp, h.log)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
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

	var userInput signInput

	if err := json.Unmarshal(body, userInput); err != nil {
		h.log.Errorf("Error failed to unmarshal request body: %v", err.Error())
		commonHttp.ErrorResponse(w, "failed to unmarshal request body", http.StatusBadRequest, h.log)
		return
	}

	token, err := h.au.LoginUser(userInput.Username, userInput.Password)
	if err != nil {
		h.log.Error(err.Error())
		commonHttp.ErrorResponse(w, "Error can't login user", http.StatusBadRequest, h.log)
		return
	}

	h.log.Infof("User login with token: %s", token)

	loginResponse := &loginResponse{JWT: token}

	commonHttp.SuccessResponse(w, loginResponse, h.log)
}

/*func (h *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
	user, err := r.Context().Value(models.User)
}*/
