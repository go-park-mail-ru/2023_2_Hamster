package http

import (
	"encoding/json"
	"net/http"

	response "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions"

	quest "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/question"

	gen "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/grpc/generated"

	genQuest "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/question/delivery/grpc/generated"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
)

type Handler struct {
	// uu  user.Usecase
	// au  auth.Usecase
	su     sessions.Usecase
	auth   gen.AuthServiceClient
	client genQuest.QuestionServiceClient
	log    logger.Logger
}

func NewHandler(
	// uu user.Usecase,
	// au auth.Usecase,
	su sessions.Usecase,
	auth gen.AuthServiceClient,
	client genQuest.QuestionServiceClient,
	log logger.Logger) *Handler {
	return &Handler{
		// uu:  uu,
		// au:  au,
		su:     su,
		auth:   auth,
		client: client,
		log:    log,
	}
}

func (h *Handler) PostAnswer(w http.ResponseWriter, r *http.Request) {
	var input quest.AnswerRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {

		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	cookie, err := r.Cookie("session_id")
	if err != nil {
		response.ErrorResponse(w, http.StatusForbidden, err, "No cookie provided", h.log)
		return
	}

	// Find session in redis
	session, err := h.su.GetSessionByCookie(r.Context(), cookie.Value)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, err, "Session doesn't exist login", h.log)
		return
	}

	_, err = h.client.CreateAnswer(r.Context(), &genQuest.AnswerRequest{
		Id:     session.UserId.String(),
		Name:   input.Name,
		Rating: input.Rating,
	})
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err, "Session doesn't exist login", h.log)
	}

	response.SuccessResponse(w, http.StatusOK, response.NilBody{})
}

func (h *Handler) CheckUserAnswer(w http.ResponseWriter, r *http.Request) {
	var input quest.QuestionNameRequest

	// Decode request Body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {

		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	// Get cookie from request
	cookie, err := r.Cookie("session_id")
	if err != nil {
		response.ErrorResponse(w, http.StatusForbidden, err, "No cookie provided", h.log)
		return
	}

	// Find session in redis
	session, err := h.su.GetSessionByCookie(r.Context(), cookie.Value)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, err, "Session doesn't exist login", h.log)
		return
	}

	ans, err := h.client.CheckUserAnswer(r.Context(), &genQuest.CheckUserAnswerRequest{
		Id:           session.UserId.String(),
		QuestionName: input.QuestionName,
	})
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err, "Can't find privious ans", h.log)
	}

	resp := quest.CheckUserAnswerResponse{
		Average: ans.Average,
	}

	response.SuccessResponse(w, http.StatusOK, resp)
}

func (h *Handler) GetStat(w http.ResponseWriter, r *http.Request) {
	var input quest.QuestionNameRequest

	// Decode request Body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {

		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	stat, err := h.client.CalculateAverageRating(r.Context(), &genQuest.CalculateAverageRatingRequest{
		QuestionName: input.QuestionName,
	})
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err, "Bad", h.log)
		return
	}

	resp := quest.AverageResponse{
		Average: stat.Average,
	}

	response.SuccessResponse(w, http.StatusOK, resp)
}
