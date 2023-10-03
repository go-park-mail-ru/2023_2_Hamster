package http

import (
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user"
)

type Handler struct {
	userService user.Usecase
	logger      logger.CustomLogger
}

type balanceResponse struct {
	Balance float64 `json:"balance"`
}

type budgetPlannedResponse struct {
	BudgetPlanned float64 `json:"planned_balance"`
}

type budgetActualResponse struct {
	BudgetActual float64 `json:"planned_balance"`
}

const (
	userIdUrlParam = "userID"
)

func NewHandler(uu user.Usecase, l logger.CustomLogger) *Handler {
	return &Handler{
		userService: uu,
		logger:      l,
	}
}

// @Summary		Get Balance
// @Tags		User
// @Description	Get User balance
// @Produce		json
// @Success		200		{object}	balanceResponse	        "Show balance"
// @Failure		400		{object}	http.Error	"Client error"
// @Failure		500		{object}	http.Error	"Server error"
// @Router		/api/user/{userID}/balance [get]
func (h *Handler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := commonHttp.GetIDFromRequest(userIdUrlParam, r)

	if err != nil {
		h.logger.Infof("invalid id: %v:", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}
	balance, err := h.userService.GetUserBalance(userID)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error get balance", http.StatusBadRequest, h.logger)
	}
	balanceResponse := &balanceResponse{Balance: balance}
	commonHttp.SuccessResponse(w, balanceResponse, h.logger)

}

// @Summary		Get Planned Budget
// @Tags		User
// @Description	Get User planned budget
// @Produce		json
// @Success		200		{object}	budgetPlannedResponse	        "Show planned budget"
// @Failure		400		{object}	http.Error	"Client error"
// @Failure		500		{object}	http.Error	"Server error"
// @Router		/api/user/{userID}/plannedBudget [get]
func (h *Handler) GetPlannedBudget(w http.ResponseWriter, r *http.Request) {
	userID, err := commonHttp.GetIDFromRequest(userIdUrlParam, r)

	if err != nil {
		h.logger.Infof("invalid id: %v:", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}
	budget, err := h.userService.GetPlannedBudget(userID)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error get planned budget", http.StatusBadRequest, h.logger)
	}
	budgetResponse := &budgetPlannedResponse{BudgetPlanned: budget}
	commonHttp.SuccessResponse(w, budgetResponse, h.logger)
}

// @Summary		Get Actual Budget
// @Tags		User
// @Description	Get User actual budget
// @Produce		json
// @Success		200		{object}	budgetActualResponse	        "Show actual budget"
// @Failure		400		{object}	http.Error	"Client error"
// @Failure		500		{object}	http.Error	"Server error"
// @Router		/api/user/{userID}/actualBudget [get]
func (h *Handler) GetCurrentBudget(w http.ResponseWriter, r *http.Request) {
	userID, err := commonHttp.GetIDFromRequest(userIdUrlParam, r)

	if err != nil {
		h.logger.Infof("invalid id: %v:", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}
	budget, err := h.userService.GetCurrentBudget(userID)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error get current budget", http.StatusBadRequest, h.logger)
	}
	budgetResponse := &budgetActualResponse{BudgetActual: budget}
	commonHttp.SuccessResponse(w, budgetResponse, h.logger)
}

// func (h* Handler) GetPlannedBudget(w http.ResponseWriter, r *http.Request) {
// 	commonHttp.
// }

// func (h* Handler) ActualBudget(w http.ResponseWriter, r *http.Request) [
// 	commonHttp.
// ]
