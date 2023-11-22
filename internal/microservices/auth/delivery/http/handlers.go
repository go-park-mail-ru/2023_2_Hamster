package http

import (
	"encoding/json"
	"net/http"
	"time"

	auth "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth"

	contextutils "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/context_utils"
	response "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions"
)

const (
	userloginUrlParam = "login"
)

type Handler struct {
	au  auth.Usecase
	uu  user.Usecase
	su  sessions.Usecase
	log logger.Logger
}

func NewHandler(
	au auth.Usecase,
	uu user.Usecase,
	su sessions.Usecase,
	log logger.Logger) *Handler {
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
// @Success		201		{object}	Response[auth.SignResponse]		"User Created"
// @Failure		400		{object}	ResponseError					"Incorrect Input"
// @Failure		429		{object}	ResponseError					"Server error"
// @Router		/api/auth/signup	[post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var signUpUser auth.SignUpInput

	// Unmarshal r.Body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&signUpUser); err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Error(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	// Sanitaize
	if err := signUpUser.CheckValid(); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err, response.InvalidBodyRequest, h.log)
		return
	}

	// Creating user
	id, login, username, err := h.au.SignUp(r.Context(), signUpUser)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("Error in sign up: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't Sign Up user", h.log)
		return
	}

	// Creating session for new user
	session, err := h.su.CreateSessionById(r.Context(), id)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("Error in sign up session creation: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't Sign Up user", h.log)
		return
	}

	// Output
	regUser := auth.SignResponse{
		ID:       session.UserId,
		Login:    login,
		Username: username,
	}

	http.SetCookie(w, response.InitCookie("session_id", session.Cookie, time.Now().Add(7*24*time.Hour), "/api"))
	response.SuccessResponse(w, http.StatusCreated, regUser)
}

// @Summary		Sign In
// @Tags			Auth
// @Description	Login account
// @Accept 		json
// @Produce		json
// @Param			userInput		body		auth.LoginInput		true		"username && password"
// @Success		202			{object}	Response[auth.SignResponse]		"User logedin"
// @Failure		400			{object}	ResponseError					"Incorrect Input"
// @Failure		429			{object}	ResponseError					"Server error"
// @Router		/api/auth/signin	[post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginUser auth.LoginInput

	// Decode request Body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginUser); err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Error(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	// Sanitaize Input
	if err := loginUser.CheckValid(); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err, response.InvalidBodyRequest, h.log)
		return
	}

	// Login user
	id, login, username, err := h.au.Login(r.Context(), loginUser.Login, loginUser.PlaintPassword)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("Error in login: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't Login user", h.log)
		return
	}

	// Registrate session in redis
	session, err := h.su.CreateSessionById(r.Context(), id)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("Error in login: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't Login user", h.log)
		return
	}

	// Create Response
	regUser := auth.SignResponse{
		ID:       id,
		Login:    login,
		Username: username,
	}

	// Set Cookie
	http.SetCookie(w, response.InitCookie("session_id", session.Cookie, time.Now().Add(7*24*time.Hour), "/api"))
	// Send http Response
	response.SuccessResponse(w, http.StatusAccepted, regUser)
}

// @Summary		Validate Auth
// @Tags			Auth
// @Description	Validate auth
// @Accept 		json
// @Produce		json
// @Param			user		body		models.User		true		"user info"
// @Success		200		{object}	Response[auth.SignResponse]	"User status"
// @Failure		401		{object}	ResponseError				"Session doesn't exist"
// @Failure		403		{object}	ResponseError				"Invalid cookie"
// @Failure		500		{object}	ResponseError				"Server error: cookie read fail"
// @Router		/api/auth/checkAuth	[post]
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Get cookie from request
	cookie, err := r.Cookie("session_id")
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("Auth check error: %v", err)
		response.ErrorResponse(w, http.StatusForbidden, err, "No cookie provided", h.log)
		return
	}

	// Find session in redis
	session, err := h.su.GetSessionByCookie(r.Context(), cookie.Value)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("Auth check error: %v", err)
		response.ErrorResponse(w, http.StatusUnauthorized, err, "Session doesn't exist login", h.log)
		return
	}

	// Get User by its Id
	user, err := h.au.GetByID(r.Context(), session.UserId)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("Auth check error: %v", err)
		response.ErrorResponse(w, http.StatusUnauthorized, err, "Can't get you username", h.log)
		return
	}

	// Create Response
	resp := auth.SignResponse{
		ID:       session.UserId,
		Username: user.Username,
	}

	// Send Response
	response.SuccessResponse(w, http.StatusOK, resp)
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
	// Get Cookie from request
	session, err := r.Cookie("session_id")
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("Log out error: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, err, "No cookie provided", h.log)
		return
	}

	// Delete key Values pair in redis
	err = h.su.DeleteSessionByCookie(r.Context(), session.Name)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("Error session delete: %v", err)
		response.ErrorResponse(w, http.StatusInternalServerError, err, "Can't delete session", h.log)
		return
	}

	// Set Http cookie
	http.SetCookie(w, response.InitCookie("session_id", "", time.Now().AddDate(0, 0, -1), "/api"))
	// Send response
	response.SuccessResponse(w, http.StatusOK, response.NilBody{})
}

// @Summary		Get unique login info
// @Tags			Auth
// @Description	Get bool parametrs about unique login
// @Produce		json
// @Success		200		{object}		Response[bool] 	"Show user"
// @Failure		400		{object}		ResponseError		"Client error"
// @Failure		500		{object}		ResponseError		"Server error"
// @Router		/api/auth/checkLogin/ [post]
func (h *Handler) CheckLoginUnique(w http.ResponseWriter, r *http.Request) {
	var userLogin auth.UniqCheckInput

	// Decode request Body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userLogin); err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Error(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	// request login
	isUnique, err := h.au.CheckLoginUnique(r.Context(), userLogin.Login)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("DB interal error")
		response.ErrorResponse(w, http.StatusInternalServerError, err, "Can't query DB", h.log)
		return
	}

	response.SuccessResponse(w, http.StatusOK, !isUnique)
}
