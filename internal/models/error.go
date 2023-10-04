package models

import (
	"fmt"

	"github.com/google/uuid"
)

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
	return fmt.Sprintf("balance from user: #%d doesn't exist", e.UserID)
}

func (e *NoSuchPlannedBudgetError) Error() string {
	return fmt.Sprintf("planned budget from user: #%d doesn't exist", e.UserID)
}

func (e *NoSuchCurrentBudget) Error() string {
	return fmt.Sprintf("actual budget from user: #%d doesn't exist", e.UserID)
}

func (e *NoSuchAccounts) Error() string {
	return fmt.Sprintf("No Such Accounts from user: #%d doesn't exist", e.UserID)
}
