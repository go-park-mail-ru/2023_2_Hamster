package goal

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Useace interface {
	CreateGoal(ctx context.Context, goal GoalCreateRequest) (uuid.UUID, error)
	UpdateGoal(ctx context.Context, goal *models.Goal) error
	DeleteGoal(ctx context.Context, goalId uuid.UUID, userId uuid.UUID) error
	GetGoals(ctx context.Context, userId uuid.UUID) ([]models.Goal, error)

	CheckGoalsState(ctx context.Context, userId uuid.UUID) ([]models.Goal, error) // check if any goals are completed
}

type Repository interface {
	CreateGoal(ctx context.Context, goal models.Goal) (uuid.UUID, error)
	UpdateGoal(ctx context.Context, goal *models.Goal) error
	DeleteGoal(ctx context.Context, goalId uuid.UUID) error
	GetGoals(ctx context.Context, userId uuid.UUID) ([]models.Goal, error)

	CheckGoalsState(ctx context.Context, userId uuid.UUID) ([]models.Goal, error) // check if any goals are completed
}
