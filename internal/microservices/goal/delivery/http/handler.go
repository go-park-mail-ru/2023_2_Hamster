package http

import (
	"encoding/json"
	"net/http"

	contextutils "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/context_utils"
	response "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/goal"
)

type Handler struct {
	goalUsecase goal.Useace
	log         logger.Logger
}

func NewHandler(gu goal.Useace, l logger.Logger) *Handler {
	return &Handler{
		goalUsecase: gu,
		log:         l,
	}
}

func (h *Handler) CreateGoal(w http.ResponseWriter, r *http.Request) {
	// get user from context
	user, err := response.GetUserFromRequest(r)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, err, response.ErrUnauthorized.Error(), h.log)
		return
	}

	// unmarshal request body
	var goalInput goal.GoalCreateRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&goalInput); err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error Corupted request body: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	goalInput.UserId = user.ID

	goalId, err := h.goalUsecase.CreateGoal(r.Context(), goalInput)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err, "can't create goal", h.log)
		return
	}

	response.SuccessResponse(w, http.StatusOK, goalId)
}

func (h *Handler) UpdateGoal(w http.ResponseWriter, r *http.Request) {
	user, err := response.GetUserFromRequest(r)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, err, response.ErrUnauthorized.Error(), h.log)
		return
	}

	var goalInput models.Goal

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&goalInput); err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error Corupted request body: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	goalInput.UserId = user.ID

	if err := h.goalUsecase.UpdateGoal(r.Context(), &goalInput); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err, "can't update goal", h.log)
		return
	}

	response.SuccessResponse(w, http.StatusOK, response.NilBody{})
}

func (h *Handler) DeleteGoal(w http.ResponseWriter, r *http.Request) {
	user, err := response.GetUserFromRequest(r)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, err, response.ErrUnauthorized.Error(), h.log)
		return
	}

	var goalInput goal.GoalDeleteRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&goalInput); err != nil {
		h.log.WithField(
			"Request-Id", contextutils.GetReqID(r.Context()),
		).Errorf("[handler] Error Corupted request body: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, err, "Corrupted request body can't unmarshal", h.log)
		return
	}
	defer r.Body.Close()

	if err := h.goalUsecase.DeleteGoal(r.Context(), goalInput.ID, user.ID); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err, "can't delete goal", h.log)
		return
	}

	response.SuccessResponse(w, http.StatusOK, response.NilBody{})
}

func (h *Handler) GetGoals(w http.ResponseWriter, r *http.Request) {
	user, err := response.GetUserFromRequest(r)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, err, response.ErrUnauthorized.Error(), h.log)
		return
	}

	goals, err := h.goalUsecase.GetGoals(r.Context(), user.ID)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err, "can't get goals", h.log)
		return
	}

	response.SuccessResponse(w, http.StatusOK, goals)
}

func (h *Handler) CheckGoalsState(w http.ResponseWriter, r *http.Request) {
	user, err := response.GetUserFromRequest(r)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, err, response.ErrUnauthorized.Error(), h.log)
		return
	}

	goals, err := h.goalUsecase.CheckGoalsState(r.Context(), user.ID)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err, "can't check goals state", h.log)
		return
	}

	response.SuccessResponse(w, http.StatusOK, goals)
}
