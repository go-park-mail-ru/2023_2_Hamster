package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/goal"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Usecase struct {
	goalRepo goal.Repository
	log      logger.Logger
}

func NewUsecase(gr goal.Repository, log logger.Logger) *Usecase {
	return &Usecase{
		goalRepo: gr,
		log:      log,
	}
}

func (u *Usecase) CreateGoal(ctx context.Context, goal models.Goal) (uuid.UUID, error) {
	goalId, err := u.goalRepo.CreateGoal(ctx, goal)
	if err != nil {
		return uuid.Nil, fmt.Errorf("[usecase] Error goal creation: %w", err)
	}

	return goalId, nil
}

func (u *Usecase) UpdateGoal(ctx context.Context, goal *models.Goal) error {
	if err := u.goalRepo.UpdateGoal(ctx, goal); err != nil {
		return fmt.Errorf("[usecase] update goal Error: %w", err)
	}

	return nil
}

func (u *Usecase) DeleteGoal(ctx context.Context, goalId uuid.UUID) error {
	if err := u.goalRepo.DeleteGoal(ctx, goalId); err != nil {
		return fmt.Errorf("[usecase] delete goal Error: %w", err)
	}

	return nil
}

func (u *Usecase) GetGoals(ctx context.Context, userId uuid.UUID) ([]models.Goal, error) {
	goals, err := u.goalRepo.GetGoals(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("[usecase] get goals Error: %w", err)
	}

	return goals, nil
}

func (u *Usecase) CheckGoalsState(ctx context.Context, userId uuid.UUID) ([]models.Goal, error) {
	goals, err := u.goalRepo.CheckGoalsState(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("[usecase] check goals state Error: %w", err)
	}

	return goals, nil
}
