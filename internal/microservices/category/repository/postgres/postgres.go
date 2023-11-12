package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

const (
	CategoryGet = `SELECT user_id, parent_tag, "name", show_income, show_outcome, regular FROM category WHERE id=$1;`

	CategoryCreate = `INSERT INTO category (user_id, parent_tag, "name", show_income, show_outcome, regular)
				      VALUES ($1, $2, $3, $4, $5, $6)
				      RETURNING id;`

	CategoryUpdate = `UPDATE category SET parent_tag=$1, "name"=$2, show_income=$3, show_outcome=$4, regular=$5 WHERE id=$6;`

	CategoryDelete = "DELETE FROM category WHERE id = $1;"

	CategoeyAll = `SELECT * FROM category WHERE user_id = $1;`

	CategoryNameCheck = `SELECT EXIST (
						SELECT "name" FROM category 
						WHERE "name" = $1 AND parent_id = $2
					);`

	CategoryExistCheck = `SELECT EXISTS (
						SELECT "name" FROM category 
						WHERE user_id = $1 AND id = $2
					);`
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

func (r *Repository) CreateTag(ctx context.Context, tag models.Category) (uuid.UUID, error) {
	row := r.db.QueryRow(ctx, CategoryCreate,
		tag.UserID,
		tag.ParentID,
		tag.Name,
		tag.ShowIncome,
		tag.ShowOutcome,
		tag.Regular,
	)
	var id uuid.UUID

	err := row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("[repo] failed create new tag: %w", err)
	}
	return id, nil
}

func (r *Repository) UpdateTag(ctx context.Context, tag *models.Category) error {
	/*var exists bool

	row := r.db.QueryRow(ctx, CategoryGet, tag.ID)
	err := row.Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("[repo] Error tag doesn't exist: %w", err)
	} else if err != nil {
		return fmt.Errorf("[repo] failed request db %s, %w", CategoryGet, err)
	} */

	_, err := r.db.Exec(ctx, CategoryUpdate,
		tag.ParentID,
		tag.Name,
		tag.ShowIncome,
		tag.ShowOutcome,
		tag.Regular,
		tag.ID,
	)
	if err != nil {
		return fmt.Errorf("[repo] failed to update category info: %s, %w", CategoryUpdate, err)
	}
	return nil
}

func (r *Repository) DeleteTag(ctx context.Context, tagId uuid.UUID) error {
	var exists bool

	row := r.db.QueryRow(ctx, CategoryGet, tagId)
	err := row.Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("[repo] tag doesn't exist Error: %v", err)
	} else if err != nil {
		return fmt.Errorf("[repo] failed request db %s, %w", CategoryGet, err)
	}

	_, err = r.db.Exec(ctx, CategoryDelete, tagId)
	if err != nil {
		return fmt.Errorf("[repo] failed to delete category %s, %w", CategoryDelete, err)
	}
	return nil
}

func (r *Repository) GetTags(ctx context.Context, userID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category

	rows, err := r.db.Query(ctx, CategoeyAll, userID)
	if err != nil {
		return nil, fmt.Errorf("[repo] Error no tags found: %v", err)
	}

	for rows.Next() {
		var tag models.Category
		if err := rows.Scan(
			&tag.ID,
			&tag.UserID,
			&tag.ParentID,
			&tag.Name,
			&tag.ShowIncome,
			&tag.ShowOutcome,
			&tag.Regular,
		); err != nil {
			return nil, err
		}

		categories = append(categories, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[repo] Error rows error: %v", err)
	}

	if len(categories) == 0 {
		return nil, fmt.Errorf("[repo] Error no tags found: %v", err)
	}
	return categories, nil
}

func (r *Repository) CheckNameUniq(ctx context.Context, userId uuid.UUID, parentId uuid.UUID, name string) (bool, error) {
	var exist bool
	err := r.db.QueryRow(ctx, CategoryNameCheck, userId, parentId, name).Scan(&exist)
	if errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("[repo] No rows found: %v", err)
	} else if err != nil {
		return false, fmt.Errorf("[repo] Error: %v", err)
	}

	if exist {
		return false, nil
	}
	return true, nil
}

func (r *Repository) CheckExist(ctx context.Context, userId uuid.UUID, tagId uuid.UUID) (bool, error) {
	var exist bool
	err := r.db.QueryRow(ctx, CategoryExistCheck, userId, tagId).Scan(&exist)
	if errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("[repo] No rows found: %v", err)
	} else if err != nil {
		return false, fmt.Errorf("[repo] Error: %v", err)
	}
	return true, nil
}
