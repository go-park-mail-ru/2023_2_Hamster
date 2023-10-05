package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth"
	"github.com/google/uuid"
)

type signUpResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
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

// @Summary		Sign Up
// @Tags			Auth
// @Description	Create Account
// @Accept 		json
// @Produce		json
// @Param			user		body		models.User		true		"user info"
// @Success		200		{object}	signUpResponse				"User Created"
// @Failure		400		{object}	http.Error				"Incorrect Input"
// @Failure		500		{object}	http.Error				"Server error"
// @Router		/api/auth/signup	[post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		h.log.Error(err.Error())
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err.Error(), h.log)
		return
	}

	h.log.Debug("request body successfully decoded\n", r)

	id, token, err := h.au.SignUpUser(user)
	if err != nil {
		h.log.Error(err.Error())
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err.Error(), h.log)
		return
	}

	h.log.Infof("User created with id: %d", id)

	suResp := signUpResponse{ID: id, Username: user.Username}

	http.SetCookie(w, &http.Cookie{
		Name:     "Authentication",
		Value:    token.Value,
		Expires:  token.Expires,
		Path:     "/",
		HttpOnly: true,
	})
	commonHttp.JSON(w, http.StatusOK, suResp)
}

// @Summary		Sign In
// @Tags			Auth
// @Description	Login account
// @Accept 		json
// @Produce		json
// @Param			userInput		body		signInput		true		"username && password"
// @Success		200			{object}	signUpResponse			"User logedin"
// @Failure		400			{object}	http.Error			"Incorrect Input"
// @Failure		500			{object}	http.Error			"Server error"
// @Router		/api/auth/signin	[post]
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var userInput signInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		h.log.Error(err.Error())
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err.Error(), h.log)
		return
	}

	h.log.Debug("request body successfully decoded", r)

	id, token, err := h.au.SignInUser(userInput.Username, userInput.Password)
	if err != nil {
		h.log.Error(err.Error())
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err.Error(), h.log)
		return
	}

	h.log.Infof("User login with token: %s", token)

	// loginResponse := &loginResponse{JWT: token.Value}

	siResp := signUpResponse{ID: id, Username: userInput.Username}

	http.SetCookie(w, &http.Cookie{
		Name:     "Authentication",
		Value:    token.Value,
		Expires:  token.Expires,
		Path:     "",
		HttpOnly: true,
	})

	commonHttp.JSON(w, http.StatusOK, siResp)
}

// @Summary		Validate Auth
// @Tags			Auth
// @Description	Validate auth
// @Accept 		json
// @Produce		json
// @Param			user		body		models.User		true		"user info"
// @Success		200		{object}	http.Error				"User status"
// @Failure		400		{object}	http.Error				"Invalid cookie"
// @Failure		500		{object}	http.Error				"Server error: cookie read fail"
// @Router		/api/auth/checkAuth	[post]
func (h *Handler) AccessVerification(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("Authentication")
	if errors.Is(err, http.ErrNoCookie) {
		h.log.Errorf("Error cookie token not found: %v", err)
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err.Error(), h.log)
		return
	} else if err != nil {
		h.log.Errorf("Error fail to get cookie token: %v", err)
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err.Error(), h.log)
		return
	}

	tim := tokenCookie.Expires

	if tim.Compare(time.Now()) <= 0 {
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err.Error(), h.log)
		return
	}

	id, username, err := h.au.ValidateAccessToken(tokenCookie.Value)
	if err != nil {
		h.log.Errorf("Error invalid jwt token: %v", err)
		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err.Error(), h.log)
		return
	}

	h.log.Info("User id: ", id)

	response := signUpResponse{
		Username: username,
		ID:       id,
	}

	commonHttp.JSON(w, http.StatusOK, response)
}

func (h *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("Authentication")
	http.SetCookie(w, &http.Cookie{
		Name:    "Authentication",
		Value:   c.Value,
		Expires: time.Now().UTC().AddDate(0, 0, -99999999),
		MaxAge:  -1,
	})
	h.log.Info("logout")
	commonHttp.JSON(w, http.StatusOK, "user loged out")
}
