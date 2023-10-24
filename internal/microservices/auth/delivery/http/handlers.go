package http

import (
	"encoding/json"
	"net/http"
	"time"

	response "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions"
	"github.com/google/uuid"
)

type RegistredUser struct {
	ID       uuid.UUID `json:"user_id"`
	username string    `json:"username"`
}

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
// @Success		200		{object}	Response[RegistredUser]			"User Created"
// @Failure		400		{object}	ResponseError					"Incorrect Input"
// @Failure		429		{object}	ResponseError					"Server error"
// @Router		/api/auth/signup	[post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		h.log.Error(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}

	user, err := h.au.SignUp(r.Context(), user)
	if err != nil {
		h.log.Errorf("Error in sign up: %v", err)
		response.ErrorResponse(w, http.StatusTooManyRequests, err, "Can't Sign Up user", h.log)
	}

	session, err := h.su.CreateSessionById(r.Context(), user.ID)
	regUser := RegistredUser{ID: session.UserId, username: user.Username}

	http.SetCookie(w, response.InitCookie(response.AuthTag, session.Cookie, time.Now().Add(7*24*time.Hour), "/api"))
	response.SuccessResponse[RegistredUser](w, http.StatusAccepted, regUser)
}

// @Summary		Sign In
// @Tags			Auth
// @Description	Login account
// @Accept 		json
// @Produce		json
// @Param			userInput		body		signInput		true		"username && password"
// @Success		200			{object}	Response[RegistredUser]			"User logedin"
// @Failure		400			{object}	ResponseError			"Incorrect Input"
// @Failure		500			{object}	ResponseError			"Server error"
// @Router		/api/auth/signin	[post]
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var user models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		h.log.Error(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
}

// @Summary		Validate Auth
// @Tags			Auth
// @Description	Validate auth
// @Accept 		json
// @Produce		json
// @Param			user		body		models.User		true		"user info"
// @Success		200		{object}	ResponseError				"User status"
// @Failure		400		{object}	ResponseError				"Invalid cookie"
// @Failure		500		{object}	ResponseError				"Server error: cookie read fail"
// @Router		/api/auth/checkAuth	[post]
func (h *Handler) AccessVerification(w http.ResponseWriter, r *http.Request) {
}

// @Summary		Validate Auth
// @Tags			Auth
// @Description	Validate auth
// @Accept 		json
// @Produce		json
// @Param			user		body		models.User		true		"user info"
// @Success		200		{object}	ResponseError				"User status"
// @Failure		400		{object}	ResponseError				"Invalid cookie"
// @Failure		500		{object}	ResponseError				"Server error: cookie read fail"
// @Router		/api/auth/checkAuth	[post]
func (h *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
}

// @Summary		Sign Up
// @Tags			Auth
// @Description	Create Account
// @Accept 		json
// @Produce		json
// @Param			user		body		models.User		true		"user info"
// @Success		200		{object}	signUpResponse				"User Created"
// @Failure		400		{object}	ResponseError				"Incorrect Input"
// @Failure		500		{object}	ResponseError				"Server error"
// @Router		/api/auth/signup	[post]
//func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
//	var user models.User
//
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&user); err != nil {
//		h.log.Error(err.Error())
//		commonHttp.ErrorResponse(w, http.StatusBadRequest, err.Error(), h.log)
//		return
//	}
//
//	h.log.Debug("request body successfully decoded\n", r)
//
//	id, err := h.au.SignUp(ctx.TODO(), user)
//	if err != nil {
//		h.log.Error(err.Error())
//		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err.Error(), h.log)
//		return
//	}
//
//	h.log.Infof("User created with id: %d", id)
//
//	suResp := signUpResponse{ID: id, Username: user.Username}
//
//	http.SetCookie(w, InitCookie(AuthCookie, token.Value, token.Expires))
//	commonHttp.JSON(w, http.StatusOK, suResp)
//}
//
//// @Summary		Sign In
//// @Tags			Auth
//// @Description	Login account
//// @Accept 		json
//// @Produce		json
//// @Param			userInput		body		signInput		true		"username && password"
//// @Success		200			{object}	signUpResponse			"User logedin"
//// @Failure		400			{object}	ResponseError			"Incorrect Input"
//// @Failure		500			{object}	ResponseError			"Server error"
//// @Router		/api/auth/signin	[post]
//func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
//	userInput := &models.SignInput{}
//
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&userInput); err != nil {
//		h.log.Error(err.Error())
//		commonHttp.ErrorResponse(w, http.StatusBadRequest, err.Error(), h.log)
//		return
//	}
//
//	h.log.Debug("request body successfully decoded", r)
//
//	id, token, err := h.au.SignInUser(userInput.Username, userInput.Password)
//	if err != nil {
//		h.log.Error(err.Error())
//		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err.Error(), h.log)
//		return
//	}
//
//	h.log.Infof("User login with token: %s", token)
//
//	// loginResponse := &loginResponse{JWT: token.Value}
//
//	siResp := signUpResponse{ID: id, Username: userInput.Username}
//
//	http.SetCookie(w, InitCookie(AuthCookie, token.Value, token.Expires))
//
//	commonHttp.JSON(w, http.StatusOK, siResp)
//}
//
//// @Summary		Validate Auth
//// @Tags			Auth
//// @Description	Validate auth
//// @Accept 		json
//// @Produce		json
//// @Param			user		body		models.User		true		"user info"
//// @Success		200		{object}	ResponseError				"User status"
//// @Failure		400		{object}	ResponseError				"Invalid cookie"
//// @Failure		500		{object}	ResponseError				"Server error: cookie read fail"
//// @Router		/api/auth/checkAuth	[post]
//func (h *Handler) AccessVerification(w http.ResponseWriter, r *http.Request) {
//	tokenCookie, err := r.Cookie(AuthCookie)
//	if errors.Is(err, http.ErrNoCookie) {
//		h.log.Errorf("Error cookie token not found: %v", err)
//		commonHttp.ErrorResponse(w, http.StatusBadRequest, err.Error(), h.log)
//		return
//	} else if err != nil {
//		h.log.Errorf("Error fail to get cookie token: %v", err)
//		commonHttp.ErrorResponse(w, http.StatusBadRequest, err.Error(), h.log)
//		return
//	}
//
//	expirationTime := tokenCookie.Expires
//	if expirationTime.IsZero() {
//		h.log.Info("Cookie does not have an expiration time.")
//	} else if expirationTime.UTC().Before(time.Now()) {
//		h.log.Errorf("Cookie has expired")
//		commonHttp.ErrorResponse(w, http.StatusUnauthorized, "Cookie has expired", h.log)
//		return
//	}
//
//	id, username, err := h.au.ValidateAccessToken(tokenCookie.Value)
//	if err != nil {
//		h.log.Errorf("Error invalid jwt token: %v", err)
//		commonHttp.ErrorResponse(w, http.StatusUnauthorized, err.Error(), h.log)
//		return
//	}
//
//	h.log.Info("User id: ", id)
//
//	response := signUpResponse{Username: username, ID: id}
//
//	commonHttp.JSON(w, http.StatusOK, response)
//}
//
//func (h *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
//	http.SetCookie(w, &http.Cookie{
//		Name:  AuthCookie,
//		Value: "",
//		Path:  "/",
//	})
//	h.log.Info("logout")
//	commonHttp.JSON(w, http.StatusOK, "user loged out")
//}
