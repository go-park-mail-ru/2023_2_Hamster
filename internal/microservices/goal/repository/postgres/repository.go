package postgres

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

const (
	GoalGet = `SELECT user_id, "name", "description", target, "date" FROM goal WHERE id=$1;`

	GoalCreate = `INSERT INTO goal (user_id, "name", "description", target, "date")
				      VALUES ($1, $2, $3, $4, $5)
				      RETURNING id;`

	GoalUpdate = `UPDATE goal SET "name"=$1, "description"=$2, target=$3, "date"=$4 WHERE id=$5;`

	GoalDelete = "DELETE FROM goal WHERE id = $1;"

	GoalAll = `SELECT * FROM goal WHERE user_id = $1;`

	GoalAllDone = `SELECT * FROM goal WHERE user_id = $1 AND WHERE state='done';`
)

type Repository struct {
	db  postgresql.DbConn
	log logger.Logger
}

func NewRepository(db postgresql.DbConn, log logger.Logger) *Repository {
	return &Repository{
		db:  db,
		log: log,
	}
}

func (r *Repository) CreateGoal(ctx context.Context, goal models.Goal) (uuid.UUID, error) {
	row := r.db.QueryRow(ctx, GoalCreate,
		goal.UserId,
		goal.Name,
		goal.Description,
		goal.Target,
		goal.Date,
	)

	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *Repository) UpdateGoal(ctx context.Context, goal *models.Goal) error {
	_, err := r.db.Exec(ctx, GoalUpdate,
		goal.Name,
		goal.Description,
		goal.Target,
		goal.Date,
		goal.ID,
	)
	if err != nil {
		return fmt.Errorf("[repo] UpdateGoal: %w", err)
	}

	return nil
}

func (r *Repository) DeleteGoal(ctx context.Context, goalId uuid.UUID) error {
	_, err := r.db.Exec(ctx, GoalDelete, goalId)
	if err != nil {
		return fmt.Errorf("[repo] Delete Goal: %w", err)
	}

	return nil
}

func (r *Repository) GetGoals(ctx context.Context, userId uuid.UUID) ([]models.Goal, error) {
	var goals []models.Goal

	rows, err := r.db.Query(ctx, GoalAll, userId)
	if err != nil {
		return nil, fmt.Errorf("[repo] Error no tags found: %v", err)
	}

	for rows.Next() {
		var goal models.Goal
		if err := rows.Scan(
			&goal.ID,
			&goal.UserId,
			&goal.Name,
			&goal.Description,
			&goal.Target,
			&goal.Date,
		); err != nil {
			return nil, err
		}

		goals = append(goals, goal)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[repo] Error rows error: %v", err)
	}

	if len(goals) == 0 {
		return nil, fmt.Errorf("[repo] Error no tags found: %v", err)
	}
	return goals, nil
}

func (r *Repository) CheckGoalsState(ctx context.Context, userId uuid.UUID) ([]models.Goal, error) {
	var goals []models.Goal

	rows, err := r.db.Query(ctx, GoalAllDone, userId)
	if err != nil {
		return nil, fmt.Errorf("[repo] Error no tags found: %v", err)
	}

	for rows.Next() {
		var goal models.Goal
		if err := rows.Scan(
			&goal.ID,
			&goal.UserId,
			&goal.Name,
			&goal.Description,
			&goal.Target,
			&goal.Date,
		); err != nil {
			return nil, err
		}

		goals = append(goals, goal)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[repo] Error rows error: %v", err)
	}

	if len(goals) == 0 {
		return nil, fmt.Errorf("[repo] Error no tags found: %v", err)
	}
	return goals, nil
}
