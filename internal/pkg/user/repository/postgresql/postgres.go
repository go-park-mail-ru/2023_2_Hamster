package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRep struct {
	db     *sqlx.DB
	logger logger.CustomLogger
}

func NewRepository(db *sqlx.DB, l logger.CustomLogger) *UserRep {
	return &UserRep{
		db:     db,
		logger: l,
	}
}

func (r *UserRep) CreateUser(u models.User) (uuid.UUID, error) {
	query := `INSERT INTO users
			 (username, password_hash, first_name, last_name, planned_budget, avatar_url)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

	row := r.db.QueryRow(query, u.Username, u.FirstName, u.LastName, u.PlannedBudget, u.Password, u.AvatarURL)
	var id uuid.UUID

	err := row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("error request %w", err)
	}
	return id, nil
}

func (r *UserRep) GetByID(userID uuid.UUID) (*models.User, error) {
	query := `SELECT id, 
				username,
				password_hash,
				first_name,
				last_name,
				planned_budget,
				avatar_url
			 FROM users
			 WHERE id = $1;`

	row := r.db.QueryRow(query, userID)
	var u models.User

	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.FirstName, &u.LastName, &u.PlannedBudget, &u.AvatarURL)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("nothing found for this request %w", err)
	} else if err != nil {
		return &models.User{},
			fmt.Errorf("failed request db %s, %w", query, err)

	}
	return &u, nil
}

func (r *UserRep) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT username,
				password_hash,
				first_name,
				last_name,
				planned_budget,
				avatar_url
			 From users WHERE (username=&1)`
	row := r.db.QueryRow(query, username)
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.FirstName, &u.LastName, &u.PlannedBudget, &u.AvatarURL)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("nothing found for this request %w", err)
	} else if err != nil {
		return &models.User{},
			fmt.Errorf("failed request db %s, %w", query, err)
	}
	return &u, nil
}
