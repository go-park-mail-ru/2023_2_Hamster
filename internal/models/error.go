package models

import (
	"fmt"

	"github.com/google/uuid"
)

type UnathorizedError struct{}

func (e *UnathorizedError) Error() string {
	return "unathorized"
}

type NoSuchUserError struct {
	UserID uuid.UUID
}

type NoSuchUserIdBalanceError struct {
	UserID uuid.UUID
}

type NoSuchCurrentBudget struct {
	UserID uuid.UUID
}

type NoSuchPlannedBudgetError struct {
	UserID uuid.UUID
}

type NoSuchAccounts struct {
	UserID uuid.UUID
}

func (e *NoSuchUserIdBalanceError) Error() string {
	return fmt.Sprintf("balance from user: %s doesn't exist", e.UserID.String())
}

func (e *NoSuchPlannedBudgetError) Error() string {
	return fmt.Sprintf("planned budget from user: %s doesn't exist", e.UserID.String())
}

func (e *NoSuchCurrentBudget) Error() string {
	return fmt.Sprintf("actual budget from user: %s doesn't exist", e.UserID.String())
}

func (e *NoSuchAccounts) Error() string {
	return fmt.Sprintf("No Such Accounts from user: %s doesn't exist", e.UserID.String())
}

func (e *NoSuchUserError) Error() string {
	return fmt.Sprintf("No Such user: %s doesn't exist", e.UserID.String())
}
