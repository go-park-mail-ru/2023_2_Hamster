package http

import (
	"errors"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
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
	BudgetActual float64 `json:"actual_balance"`
}

type account struct {
	Account []models.Accounts
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
// @Tags			User
// @Description	Get User balance
// @Produce		json
// @Success		200		{object}	http.Response  "Show balance"
// @Failure		400		{object}	http.Error	"Client error"
// @Failure		500		{object}	http.Error	"Server error"
// @Router		/api/user/{userID}/balance [get]
func (h *Handler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := commonHttp.GetIDFromRequest(userIdUrlParam, r)

	if err != nil {
		h.logger.Infof("invalid id: %v:", err)

		commonHttp.ErrorResponse(w, http.StatusBadRequest, err.Error(), h.logger)
		return
	}
	balance, err := h.userService.GetUserBalance(userID)

	var errNoSuchUserIdBalanceError *models.NoSuchUserIdBalanceError
	if errors.As(err, &errNoSuchUserIdBalanceError) {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	commonHttp.JSON(w, http.StatusOK, balanceResponse{Balance: balance})
}

// @Summary		Get Planned Budget
// @Tags			User
// @Description	Get User planned budget
// @Produce		json
// @Success		200		{object}	budgetPlannedResponse	"Show planned budget"
// @Failure		400		{object}	http.Error			"Client error"
// @Failure		500		{object}	http.Error			"Server error"
// @Router		/api/user/{userID}/plannedBudget [get]
func (h *Handler) GetPlannedBudget(w http.ResponseWriter, r *http.Request) {
	userID, err := commonHttp.GetIDFromRequest(userIdUrlParam, r)

	if err != nil {
		h.logger.Infof("invalid id: %v:", err)

		commonHttp.ErrorResponse(w, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	budget, err := h.userService.GetPlannedBudget(userID)

	var errNoSuchPlannedBudgetError *models.NoSuchPlannedBudgetError
	if errors.As(err, &errNoSuchPlannedBudgetError) {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	commonHttp.JSON(w, http.StatusOK, budgetPlannedResponse{BudgetPlanned: budget})

}

// @Summary		Get Actual Budget
// @Tags			User
// @Description	Get User actual budget
// @Produce		json
// @Success		200		{object}	budgetActualResponse	"Show actual budget"
// @Failure		400		{object}	http.Error			"Client error"
// @Failure		500		{object}	http.Error			"Server error"
// @Router		/api/user/{userID}/actualBudget [get]
func (h *Handler) GetCurrentBudget(w http.ResponseWriter, r *http.Request) {
	userID, err := commonHttp.GetIDFromRequest(userIdUrlParam, r)

	if err != nil {
		h.logger.Infof("invalid id: %v:", err)
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	budget, err := h.userService.GetCurrentBudget(userID)

	var errNoSuchCurrentBudget *models.NoSuchCurrentBudget
	if errors.As(err, &errNoSuchCurrentBudget) {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	commonHttp.JSON(w, http.StatusOK, budgetActualResponse{BudgetActual: budget})
}

// @Summary		Get User Accounts
// @Tags			User
// @Description	Get User accounts
// @Produce		json
// @Success		200		{object}	account	     	"Show actual accounts"
// @Failure		400		{object}	http.Error		"Client error"
// @Failure		500		{object}	http.Error		"Server error"
// @Router		/api/user/{userID}/accounts/all [get]
func (h *Handler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	userID, err := commonHttp.GetIDFromRequest(userIdUrlParam, r)

	if err != nil {
		h.logger.Infof("invalid id: %v:", err)

		commonHttp.ErrorResponse(w, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	accountInfo, err := h.userService.GetAccounts(userID)

	var errNoSuchAccounts *models.NoSuchAccounts

	if errors.As(err, &errNoSuchAccounts) {
		h.logger.Error(err.Error())
		commonHttp.JSON(w, http.StatusOK, "")
		return
	}

	if err != nil {
		commonHttp.ErrorResponse(w, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	budgetResponse := &account{Account: accountInfo}
	commonHttp.JSON(w, http.StatusOK, budgetResponse)
}
