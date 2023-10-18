package transfer_models

import (
	"html"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

// Response error message
const (
	BalanceNotFound       = "no such balance"
	PlannedBudgetNotFound = "no such planned budget"
	CurrentBudgetNotFound = "no such current budget"
	AccountNotFound       = "no such account"
	UserNotFound          = "no such user"
	UserFeedNotFound      = "no such feed info"

	BalanceGetServerError       = "can't get balance"
	PlannedBudgetGetServerError = "can't get planned budget"
	CurrentBudgetGetServerError = "can't get current budget"
	AccountServerError          = "can't get account"
	UserServerError             = "can't get user"
	UserFeedServerError         = "can't get feed info"
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
}

type UserTransfer struct {
	ID            uuid.UUID `json:"id" valid:""`
	Username      string    `json:"username" valid:"required,maxstringlength(20)"`
	PlannedBudget float64   `json:"planned_budget" valid:"required,float"`
	AvatarURL     string    `json:"avatar_url" valid:""`
}

func (ui *UserTransfer) CheckValid() error {
	ui.Username = html.EscapeString(ui.Username)
	_, err := valid.ValidateStruct(*ui)

	return err
}

func (ui *UserTransfer) ToUser(user *models.User) *models.User {
	return &models.User{
		ID:            user.ID,
		Username:      ui.Username,
		PlannedBudget: ui.PlannedBudget,
		Password:      user.Password,
		AvatarURL:     user.AvatarURL,
		Salt:          user.Salt,
	}
}

func InitUserTransfer(user models.User) UserTransfer {
	return UserTransfer{
		ID:            user.ID,
		Username:      user.Username,
		PlannedBudget: user.PlannedBudget,
		AvatarURL:     user.AvatarURL,
	}
}
