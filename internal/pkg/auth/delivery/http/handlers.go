package http

import (
	"encoding/json"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth"
	"github.com/google/uuid"
)

type signUpResponse struct {
	ID uuid.UUID `json:"id"`
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

func NewHandler(au auth.Usecase, log logger.CustomLogger) *Handler {
	return &Handler{
		au:  au,
		log: log,
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		h.log.Error(err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.log)
		return
	}

	h.log.Debug("request body successfully decoded\n", r)

	id, err := h.au.SignUpUser(user)
	if err != nil {
		h.log.Error(err.Error())
		commonHttp.ErrorResponse(w, "server error", http.StatusBadRequest, h.log)
		return
	}

	h.log.Infof("User created with id: %d", id)

	suResp := signUpResponse{ID: id}

	commonHttp.SuccessResponse(w, suResp, h.log)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var userInput signInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		h.log.Error(err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.log)
		return
	}

	h.log.Debug("request body successfully decoded", r)

	token, err := h.au.SignInUser(userInput.Username, userInput.Password)
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
