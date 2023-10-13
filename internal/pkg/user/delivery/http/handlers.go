package http

import (
	"errors"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/http"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/delivery/http/transfer_models"
)

type Handler struct {
	userService user.Usecase
	logger      logger.CustomLogger
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
// @Success		200		{object}	Response[transfer_models.BalanceResponse] "Show balance"
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

	response := transfer_models.BalanceResponse{Balance: balance}
	commonHttp.JSON[transfer_models.BalanceResponse](w, http.StatusOK, response)
}

// @Summary		Get Planned Budget
// @Tags			User
// @Description	Get User planned budget
// @Produce		json
// @Success		200		{object} 	Response[transfer_models.BudgetPlannedResponse]	"Show planned budget"
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

	response := transfer_models.BudgetPlannedResponse{BudgetPlanned: budget}
	commonHttp.JSON[transfer_models.BudgetPlannedResponse](w, http.StatusOK, response)
}

// @Summary		Get Actual Budget
// @Tags			User
// @Description	Get User actual budget
// @Produce		json
// @Success		200		{object}	Response[transfer_models.BudgetActualResponse]	"Show actual budget"
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

	response := transfer_models.BudgetActualResponse{BudgetActual: budget}
	commonHttp.JSON[transfer_models.BudgetActualResponse](w, http.StatusOK, response)
}

// @Summary		Get User Accounts
// @Tags			User
// @Description	Get User accounts
// @Produce		json
// @Success		200		{object}	Response[transfer_models.Account]	     	"Show actual accounts"
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

	response := transfer_models.Account{Account: accountInfo}
	commonHttp.JSON[transfer_models.Account](w, http.StatusOK, response)
}

// @Summary		Get User Accounts
// @Tags			User
// @Description	Get User accounts
// @Produce		json
// @Success		200		{object}	Response[transfer_models.UserFeed]	     	"Show actual accounts"
// @Failure		400		{object}	http.Error		"Client error"
// @Failure		500		{object}	http.Error		"Server error"
// @Router		/api/user/{userID}/feed [get]
func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) { // need test
	status := http.StatusOK
	userID, err := commonHttp.GetIDFromRequest(userIdUrlParam, r)

	if err != nil {
		h.logger.Infof("invalid id: %v:", err)

		commonHttp.ErrorResponse(w, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	dataFeed, multierr := h.userService.GetFeed(userID)
	if multierr.ErrorOrNil() != nil {
		h.logger.Error(multierr)
		status = http.StatusInternalServerError
	}

	commonHttp.JSON(w, status, dataFeed)
}
