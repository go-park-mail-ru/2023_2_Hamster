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
			 (username, password_hash, salt)
		VALUES ($1, $2, $3) RETURNING id;`

	row := r.db.QueryRow(query, u.Username, u.Password, u.Salt)
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
				planned_budget,
				avatar_url
			 FROM users
			 WHERE id = $1;`

	row := r.db.QueryRow(query, userID)
	var u models.User

	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.PlannedBudget, &u.AvatarURL)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("nothing found for this request %w", err)
	} else if err != nil {
		return &models.User{},
			fmt.Errorf("failed request db %s, %w", query, err)

	}
	return &u, nil
}

func (r *UserRep) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id,
				username,
				password_hash,
				planned_budget,
				avatar_url,
				salt
			 From users WHERE (username=$1)`
	row := r.db.QueryRow(query, username)
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.PlannedBudget, &u.AvatarURL, &u.Salt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("[repository] nothing found for this request %w", err)
	} else if err != nil {
		return &models.User{},
			fmt.Errorf("[repository] failed request db %w", err)
	}
	return &u, nil
}

func (r *UserRep) GetUserBalance(userID uuid.UUID) (float64, error) {
	var totalBalance float64
	err := r.db.QueryRow("SELECT SUM(balance) FROM accounts WHERE user_id = $1", userID).Scan(&totalBalance)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("[repository] nothing found for this request %w", err)
	} else if err != nil {
		return 0, fmt.Errorf("[repository] failed request db %w", err)
	}

	return totalBalance, nil
}

func (r *UserRep) GetPlannedBudget(userID uuid.UUID) (float64, error) {
	var plannedBudget float64
	err := r.db.QueryRow("SELECT planned_budget FROM users WHERE id = $1", userID).Scan(&plannedBudget)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("[repository] nothing found for this request %w", err)
	} else if err != nil {
		return 0, fmt.Errorf("[repository] failed request db %w", err)
	}

	return plannedBudget, nil
}

func (r *UserRep) GetCurrentBudget(userID uuid.UUID) (float64, error) {
	var currentBudget float64

	err := r.db.QueryRow(`SELECT SUM(total) AS total_sum
					  FROM Transaction
					  WHERE date_part('month', date) = date_part('month', CURRENT_DATE)
  					  AND date_part('year', date) = date_part('year', CURRENT_DATE)
					  AND is_income = false
					  AND user_id = $1;`, userID).Scan(&currentBudget)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("[repository] nothing found for this request %w", err)
	} else if err != nil {
		return 0, fmt.Errorf("[repository] failed request db %w", err)
	}

	return currentBudget, nil
}
