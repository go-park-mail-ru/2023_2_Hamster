package http

import (
	"encoding/json"
	"net/http"
	"time"

	response "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	auth "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions"
)

const (
	userloginUrlParam = "login"
)

type Handler struct {
	au  auth.Usecase
	uu  user.Usecase
	su  sessions.Usecase
	log logger.CustomLogger
}

func NewHandler(
	au auth.Usecase,
	uu user.Usecase,
	su sessions.Usecase,
	log logger.CustomLogger) *Handler {
	return &Handler{
		au:  au,
		uu:  uu,
		su:  su,
		log: log,
	}
}

// @Summary		Sign Up
// @Tags			Auth
// @Description	Create Account
// @Accept 		json
// @Produce		json
// @Param			user		body		models.User			true		"user info"
// @Success		200		{object}	Response[auth.SignResponse]			"User Created"
// @Failure		400		{object}	ResponseError					"Incorrect Input"
// @Failure		429		{object}	ResponseError					"Server error"
// @Router		/api/auth/signup	[post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var signUpUser auth.SignUpInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&signUpUser); err != nil {
		h.log.Error(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	id, username, err := h.au.SignUp(r.Context(), signUpUser)
	if err != nil {
		h.log.Errorf("Error in sign up: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't Sign Up user", h.log)
		return
	}

	session, err := h.su.CreateSessionById(r.Context(), id)
	if err != nil {
		h.log.Errorf("Error in sign up session creation: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't Sign Up user", h.log)
		return
	}

	regUser := auth.SignResponse{
		ID:       session.UserId,
		Username: username,
	}

	http.SetCookie(w, response.InitCookie("session_id", session.Cookie, time.Now().Add(7*24*time.Hour), "/api"))
	response.SuccessResponse[auth.SignResponse](w, http.StatusAccepted, regUser)
}

// @Summary		Sign In
// @Tags			Auth
// @Description	Login account
// @Accept 		json
// @Produce		json
// @Param			userInput		body		auth.LoginInput		true		"username && password"
// @Success		200			{object}	Response[auth.SignResponse]			"User logedin"
// @Failure		400			{object}	ResponseError			"Incorrect Input"
// @Failure		500			{object}	ResponseError			"Server error"
// @Router		/api/auth/signin	[post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginUser auth.LoginInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginUser); err != nil {
		h.log.Error(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	id, login, err := h.au.Login(r.Context(), loginUser.Login, loginUser.PlaintPassword)
	if err != nil {
		h.log.Errorf("Error in login: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't Login user", h.log)
		return
	}

	session, err := h.su.CreateSessionById(r.Context(), id)
	if err != nil {
		h.log.Errorf("Error in login: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't Login user", h.log)
		return
	}

	regUser := auth.SignResponse{
		ID:       id,
		Username: login,
	}

	http.SetCookie(w, response.InitCookie("session_id", session.Cookie, time.Now().Add(7*24*time.Hour), "/api"))
	response.SuccessResponse[auth.SignResponse](w, http.StatusAccepted, regUser)
}

// @Summary		Validate Auth
// @Tags			Auth
// @Description	Validate auth
// @Accept 		json
// @Produce		json
// @Param			user		body		models.User		true		"user info"
// @Success		200		{object}	Response[auth.SignResponse]	"User status"
// @Failure		400		{object}	ResponseError				"Invalid cookie"
// @Failure		500		{object}	ResponseError				"Server error: cookie read fail"
// @Router		/api/auth/checkAuth	[post]
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		h.log.Errorf("Auth check error: %v", err)
		response.ErrorResponse(w, http.StatusForbidden, err, "No cookie provided", h.log)
		return
	}

	session, err := h.su.GetSessionByCookie(r.Context(), cookie.Value)
	if err != nil {
		h.log.Errorf("Auth check error: %v", err)
		response.ErrorResponse(w, http.StatusUnauthorized, err, "Session doesn't exist login", h.log)
		return
	}

	response.SuccessResponse[models.Session](w, http.StatusOK, session)
}

// @Summary		Validate Auth
// @Tags			Auth
// @Description	Validate auth
// @Accept 		json
// @Produce		json
// @Param			user		body		models.User		true		"user info"
// @Success		200		{object}	Response[auth.SignResponse]   "User status"
// @Failure		400		{object}	ResponseError				"Invalid cookie"
// @Failure		500		{object}	ResponseError				"Server error: cookie read fail"
// @Router		/api/auth/checkAuth	[post]
func (h *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err != nil {
		h.log.Errorf("Log out error: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, err, "No cookie provided", h.log)
		return
	}

	err = h.su.DeleteSessionByCookie(r.Context(), session.Name)
	if err != nil {
		h.log.Errorf("Error session delete: %v", err)
		response.ErrorResponse(w, http.StatusInternalServerError, err, "Can't delete session", h.log)
		return
	}

	http.SetCookie(w, response.InitCookie("session", "", time.Now().AddDate(0, 0, -1), "/api"))
	response.SuccessResponse[response.NilBody](w, http.StatusOK, response.NIL())
}

// @Summary		Get unique login info
// @Tags		Auth
// @Description	Get bool parametrs about unique login
// @Produce		json
// @Success		200		{object}	Response[bool] "Show user"
// @Failure		400		{object}	ResponseError	"Client error"
// @Failure		500		{object}	ResponseError	"Server error"
// @Router		/api/auth/check-unique-login/{login} [get]
func (h *Handler) CheckLoginUnique(w http.ResponseWriter, r *http.Request) {
	userLogin := response.GetloginFromRequest(userloginUrlParam, r)
	isUnique, err := h.au.CheckLoginUnique(r.Context(), userLogin)

	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err, "Can't get unique info login", h.log)
		return
	}
	response.SuccessResponse(w, http.StatusOK, isUnique)
}
