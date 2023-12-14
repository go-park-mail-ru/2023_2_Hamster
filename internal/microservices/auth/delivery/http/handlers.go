package http

import (
	"encoding/json"
	"net/http"
	"time"

	auth "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth"
	gen "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	contextutils "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/context_utils"
	response "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
)

type Handler struct {
	// uu  user.Usecase
	// au  auth.Usecase
	su     sessions.Usecase
	client gen.AuthServiceClient
	log    logger.Logger
}

func NewHandler(
	// uu user.Usecase,
	// au auth.Usecase,
	su sessions.Usecase,
	cl gen.AuthServiceClient,
	log logger.Logger) *Handler {
	return &Handler{
		// uu:  uu,
		// au:  au,
		su:     su,
		client: cl,
		log:    log,
	}
}

// @Summary		Sign Up
// @Tags			Auth
// @Description	Create Account
// @Accept 		json
// @Produce		json
// @Param			user		body		models.User				true	"user info"
// @Success		201		{object}	Response[auth.SignResponse]				"User Created"
// @Failure		400		{object}	ResponseError							"Incorrect Input"
// @Failure		409		{object}	ResponseError							"User already exists"
// @Failure		429		{object}	ResponseError							"Server error"
// @Router		/api/auth/signup	[post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var signUpUser auth.SignUpInput

	// Unmarshal request.Body
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
	userMeta, err := h.client.SignUp(r.Context(), &gen.SignUpRequest{
		Login:    signUpUser.Login,
		Username: signUpUser.Username,
		Password: signUpUser.PlaintPassword,
	})
	if err != nil {
		// errUserAlreadyExists := status.Error(codes.AlreadyExists, "Alredy exist") //*models.UserAlreadyExistsError
		if status.Code(err) == codes.AlreadyExists {
			response.ErrorResponse(w, http.StatusConflict, err, "User already exists", h.log)
			return
		}
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("Error in sign up: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't Sign Up user", h.log)
		return
	}
	userId, err := uuid.Parse(userMeta.Body.Id)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err, "Can't Sign Up user", h.log)
		return
	}
	// Creating session for new user
	session, err := h.su.CreateSessionById(r.Context(), userId)
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
		Login:    userMeta.Body.Login,
		Username: userMeta.Body.Username,
	}

	http.SetCookie(w, response.InitCookie("session_id", session.Cookie, time.Now().Add(7*24*time.Hour), "/api"))
	response.SuccessResponse(w, http.StatusCreated, regUser)
}

// @Summary		Sign In
// @Tags			Auth
// @Description	Login account
// @Accept 		json
// @Produce		json
// @Param		userInput		body		auth.LoginInput		true		"username && password"
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

	// Login user loginUser.Login, loginUser.PlaintPassword
	userMeta, err := h.client.Login(r.Context(), &gen.LoginRequest{
		Login:    loginUser.Login,
		Password: loginUser.PlaintPassword,
	})
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("Error in login: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't Login user", h.log)
		return
	}

	userId, err := uuid.Parse(userMeta.Body.Id)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err, "Can't Login user", h.log)
		return
	}

	// Registrate session in redis
	session, err := h.su.CreateSessionById(r.Context(), userId)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("Error in login: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't Login user", h.log)
		return
	}

	// Create Response
	regUser := auth.SignResponse{
		ID:       userId,
		Login:    userMeta.Body.Login,
		Username: userMeta.Body.Username,
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

	// Get User by its Id   session.UserId
	user, err := h.client.GetByID(r.Context(), &gen.UserIdRequest{Id: session.UserId.String()})
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
		Login:    user.Body.Login,
		Username: user.Body.Username,
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

	if err := userLogin.CheckValid(); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err, response.InvalidBodyRequest, h.log)
		return
	}

	// request login userLogin.Login
	isUnique, err := h.client.CheckLoginUnique(r.Context(), &gen.UniqCheckRequest{Login: userLogin.Login})
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("DB interal error")
		response.ErrorResponse(w, http.StatusInternalServerError, err, "Can't query DB", h.log)
		return
	}

	response.SuccessResponse(w, http.StatusOK, isUnique)
}

// @Summary		Get unique login info
// @Tags			Auth
// @Description	Get bool parametrs about unique login
// @Produce		json
// @Success		200		{object}		Response[bool] 	"Show user"
// @Failure		400		{object}		ResponseError		"Client error"
// @Failure		500		{object}		ResponseError		"Server error"
// @Router		/api/auth/checkLogin/ [post]
func (h *Handler) GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	var userId auth.UserIdInput

	// Decode request Body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userId); err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Error(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	// request login userLogin.Login
	isUnique, err := h.client.CheckLoginUnique(r.Context(), &gen.UniqCheckRequest{Login: userId.ID.String()})
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("DB interal error")
		response.ErrorResponse(w, http.StatusInternalServerError, err, "Can't query DB", h.log)
		return
	}

	response.SuccessResponse(w, http.StatusOK, isUnique)
}

// @Summary		Change Password
// @Tags		Auth
// @Description	Takes old password and newpassword and chnge password
// @Accept 		json
// @Produce		json
// @Param		userInput		body		auth.ChangePasswordInput		true		"username && password"
// @Success		200		{object}		Response[auth.SignResponse] 			"user Info"
// @Failure		400		{object}		ResponseError							"Client error"
// @Failure		500		{object}		ResponseError							"Server error"
// @Router		/api/auth/password/ [put]
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, err := response.GetUserFromRequest(r)
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error: %v", err)
		response.ErrorResponse(w, http.StatusUnauthorized, err, response.ErrUnauthorized.Error(), h.log)
		return
	}

	var changePassword auth.ChangePasswordInput

	// Decode request Body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&changePassword); err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Error(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	// Sanitaize Input
	if err := changePassword.CheckValid(); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err, response.InvalidBodyRequest, h.log)
		return
	}

	// Change Password
	_, err = h.client.ChangePassword(r.Context(), &gen.ChangePasswordRequest{
		Login:       user.Login,
		OldPassword: changePassword.OldPassword,
		NewPassword: changePassword.NewPassword,
	})
	if err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("Error in change password: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't Change Password", h.log)
		return
	}

	// Create Response
	regUser := auth.SignResponse{
		ID:       user.ID,
		Login:    user.Login,
		Username: user.Username,
	}

	// Send Response
	response.SuccessResponse(w, http.StatusOK, regUser)
}
