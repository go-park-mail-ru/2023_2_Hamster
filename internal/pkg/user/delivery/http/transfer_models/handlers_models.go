package transfer_models

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
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
	Account []models.Accounts `json:"account"`
}

type UserFeed struct {
	Account
	BalanceResponse
	BudgetPlannedResponse
	BudgetActualResponse
	ErrMes string `json:"err_message"`
}
