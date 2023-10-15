package transfer_models

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
)

// Response error message
const (
	BalanceNotFound       = "no such balance"
	PlannedBudgetNotFound = "no such planned budget"
	CurrentBudgetNotFound = "no such current budget"
	AccountNotFound       = "no such account"
	UserNotFound          = "no such user"

	BalanceGetServerError       = "can't get balance"
	PlannedBudgetGetServerError = "can't get planned budget"
	CurrentBudgetGetServerError = "can't get current budget"
	AccountServerError          = "can't get account"
	UserServerError             = "can't get user"
)

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}

type BudgetPlannedResponse struct {
	BudgetPlanned float64 `json:"planned_balance"`
}

type BudgetActualResponse struct {
	BudgetActual float64 `json:"actual_balance"`
}

type Account struct {
	AccountMas []models.Accounts `json:"account"`
}

type UserFeed struct {
	Account
	BalanceResponse
	BudgetPlannedResponse
	BudgetActualResponse
	ErrMes string `json:"err_message"`
}
