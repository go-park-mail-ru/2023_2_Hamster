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

// @Summary        Create Goal
// @Tags           Goals
// @Description    Create a new goal for the authenticated user
// @Accept         json
// @Produce        json
// @Security       ApiKeyAuth
// @Param          goalInput body goal.GoalCreateRequest true "Goal creation input"
// @Success        200 {object} Response[uuid.UUID] "Successfully created goal"
// @Failure        400 {object} ResponseError "Bad Request: Invalid request body"
// @Failure        401 {object} ResponseError "Unauthorized: Invalid or expired token"
// @Failure        500 {object} ResponseError "Internal Server Error: Failed to create goal"
// @Router         /api/user/goal/add [post]
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

// @Summary        Update Goal
// @Tags           Goals
// @Description    Update an existing goal for the authenticated user
// @Accept         json
// @Produce        json
// @Security       ApiKeyAuth
// @Param          goalInput body models.Goal true "Updated goal information"
// @Success        200 {object} Response[response.NilBody] "Successfully updated goal"
// @Failure        400 {object} ResponseError "Bad Request: Invalid request body"
// @Failure        401 {object} ResponseError "Unauthorized: Invalid or expired token"
// @Failure        404 {object} ResponseError "Not Found: Goal not found"
// @Failure        500 {object} ResponseError "Internal Server Error: Failed to update goal"
// @Router         /api/goals [put]
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

// @Summary        Delete Goal
// @Tags           Goals
// @Description    Delete an existing goal for the authenticated user
// @Accept         json
// @Produce        json
// @Security       ApiKeyAuth
// @Param          goalID path string true "ID of the goal to delete"
// @Success        200 {object} Response[response.NilBody] "Successfully deleted goal"
// @Failure        400 {object} ResponseError "Bad Request: Invalid request body"
// @Failure        401 {object} ResponseError "Unauthorized: Invalid or expired token"
// @Failure        404 {object} ResponseError "Not Found: Goal not found"
// @Failure        500 {object} ResponseError "Internal Server Error: Failed to delete goal"
// @Router         /api/user/goal/{goalID} [delete]
func (h *Handler) DeleteGoal(w http.ResponseWriter, r *http.Request) {
	user, err := response.GetUserFromRequest(r)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, err, response.ErrUnauthorized.Error(), h.log)
		return
	}

	goalID, err := response.GetIDFromRequest("goalID", r)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err, "no id provided in path", h.log)
		return
	}

	if err := h.goalUsecase.DeleteGoal(r.Context(), goalID, user.ID); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err, "can't delete goal", h.log)
		return
	}

	response.SuccessResponse(w, http.StatusOK, response.NilBody{})
}

// @Summary        Get User Goals
// @Tags           Goals
// @Description    Retrieve goals for the authenticated user
// @Accept         json
// @Produce        json
// @Security       ApiKeyAuth
// @Success        200 {object} Response[models.Goal] "Successfully retrieved user goals"
// @Failure        401 {object} ResponseError "Unauthorized: Invalid or expired token"
// @Failure        500 {object} ResponseError "Internal Server Error: Failed to get user goals"
// @Router         /api/user/goal/ [get]
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

// @Summary        Check Goals State
// @Tags           Goals
// @Description    Check the state of goals for the authenticated user
// @Accept         json
// @Produce        json
// @Security       ApiKeyAuth
// @Success        200 {object} Response[models.GoalState] "Successfully checked goals state"
// @Failure        401 {object} ResponseError "Unauthorized: Invalid or expired token"
// @Failure        400 {object} ResponseError "Bad Request: Failed to check goals state"
// @Router         /api/goals/checkState [get]
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
