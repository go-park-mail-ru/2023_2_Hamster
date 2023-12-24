package models

import (
	"fmt"

	"github.com/google/uuid"
)

type UnathorizedError struct{}

func (e *UnathorizedError) Error() string {
	return "unauthorized"
}

// =========================================UserError================================================

type DuplicateError struct {
}

type NoSuchUserInLogin struct {
	Login string
}

type NoSuchUserError struct {
	UserID uuid.UUID
}

type NoSuchTransactionError struct {
	UserID uuid.UUID
}

type NoSuchUserIdBalanceError struct {
	UserID uuid.UUID
}

// type NoSuchCurrentBudget struct {
// 	UserID uuid.UUID
// }

type NoSuchPlannedBudgetError struct {
	UserID uuid.UUID
}

type NoSuchAccounts struct {
	UserID uuid.UUID
}

type ForbiddenUserError struct{}

func (e *ForbiddenUserError) Error() string {
	return "user has no rights"
}

func (e *NoSuchUserIdBalanceError) Error() string {
	return fmt.Sprintf("balance from user: %s doesn't exist", e.UserID.String())
}

func (e *NoSuchPlannedBudgetError) Error() string {
	return fmt.Sprintf("planned budget from user: %s doesn't exist", e.UserID.String())
}

// func (e *NoSuchCurrentBudget) Error() string {
// 	return fmt.Sprintf("actual budget from user: %s doesn't exist", e.UserID.String())
// }

func (e *NoSuchAccounts) Error() string {
	return fmt.Sprintf("No Such Accounts from user: %s doesn't exist", e.UserID.String())
}

func (e *NoSuchUserError) Error() string {
	return fmt.Sprintf("No Such user: %s doesn't exist", e.UserID.String())
}

func (e *NoSuchTransactionError) Error() string {
	return fmt.Sprintf("No Such transaction: %s doesn't exist", e.UserID.String())
}

func (e *NoSuchUserInLogin) Error() string {
	return fmt.Sprintf("No Such user in login %s doesn't exist", e.Login)
}

func (e *DuplicateError) Error() string {
	return "Duplicate rows"
}

// =========================================UserError================================================

type ErrNoTags struct{}

func (e *ErrNoTags) Error() string {
	return "Don't have tags"
}
