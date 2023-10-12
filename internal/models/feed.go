package models

type UserFeed struct {
	MassAccount   []Accounts
	PlannedBudget float64
	CurrentBudget float64
	Balance       float64
}
