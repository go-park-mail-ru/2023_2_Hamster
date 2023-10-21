package postgresql

import (
	"context"
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

func (r *UserRep) CreateUser(ctx context.Context, u models.User) (uuid.UUID, error) { // need test

	query := `INSERT INTO users
			 (login, username, password_hash, salt)
		VALUES ($1, $2, $3) RETURNING id;`

	row := r.db.QueryRowContext(ctx, query, u.Login, u.Username, u.Password, u.Salt)
	var id uuid.UUID

	err := row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("error request %w", err)
	}
	return id, nil
}

func (r *UserRep) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) { // need test
	query := `SELECT *
			 FROM users
			 WHERE id = $1;`

	row := r.db.QueryRowContext(ctx, query, userID)
	var u models.User

	err := row.Scan(&u.ID, &u.Login, &u.Username, &u.Password, &u.Salt, &u.PlannedBudget, &u.AvatarURL)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("[repo] %w: %v", &models.NoSuchUserError{UserID: userID}, err)
	} else if err != nil {
		return &models.User{},
			fmt.Errorf("failed request db %s, %w", query, err)

	}
	return &u, nil
}

func (r *UserRep) GetUserByUsername(ctx context.Context, username string) (*models.User, error) { // need test
	query := `SELECT id,
				username,
				password_hash,
				planned_budget,
				avatar_url,
				salt
			 From users WHERE (username=$1)`
	row := r.db.QueryRowContext(ctx, query, username)
	var u models.User
	err := row.Scan(&u.ID, &u.Login, &u.Username, &u.Password, &u.PlannedBudget, &u.AvatarURL, &u.Salt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("[repository] nothing found for this request %w", err)
	} else if err != nil {
		return &models.User{},
			fmt.Errorf("[repository] failed request db %w", err)
	}
	return &u, nil
}

func (r *UserRep) GetUserBalance(ctx context.Context, userID uuid.UUID) (float64, error) { // need test
	var totalBalance sql.NullFloat64
	err := r.db.QueryRowContext(ctx, "SELECT SUM(balance) FROM accounts WHERE user_id = $1", userID).Scan(&totalBalance)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("[repo] %w: %v", &models.NoSuchUserIdBalanceError{UserID: userID}, err)
	} else if err != nil {
		return 0, fmt.Errorf("[repository] failed request db %w", err)
	}

	if totalBalance.Valid {
		return totalBalance.Float64, nil
	}

	return 0, nil
}

func (r *UserRep) GetPlannedBudget(ctx context.Context, userID uuid.UUID) (float64, error) { // need test
	var plannedBudget sql.NullFloat64
	err := r.db.QueryRowContext(ctx, "SELECT planned_budget FROM users WHERE id = $1", userID).Scan(&plannedBudget)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("[repo] %w: %v", &models.NoSuchPlannedBudgetError{UserID: userID}, err)
	} else if err != nil {
		return 0, fmt.Errorf("[repository] failed request db %w", err)
	}

	if plannedBudget.Valid {
		return plannedBudget.Float64, nil
	}
	return 0, nil
}

func (r *UserRep) GetCurrentBudget(ctx context.Context, userID uuid.UUID) (float64, error) { // need test
	var currentBudget sql.NullFloat64

	err := r.db.QueryRowContext(ctx, `SELECT SUM(total) AS total_sum
					  FROM Transaction
					  WHERE date_part('month', date) = date_part('month', CURRENT_DATE)
  					  AND date_part('year', date) = date_part('year', CURRENT_DATE)
					  AND is_income = false
					  AND user_id = $1;`, userID).Scan(&currentBudget)

	if err != nil {
		return 0, fmt.Errorf("[repository] failed request db %w", err)
	}

	if currentBudget.Valid {
		return currentBudget.Float64, nil
	}
	return 0, nil
}

func (r *UserRep) GetAccounts(ctx context.Context, user_id uuid.UUID) ([]models.Accounts, error) { // need test

	var accounts []models.Accounts

	rows, err := r.db.QueryContext(ctx, `SELECT * FROM accounts WHERE user_id = $1`, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account models.Accounts
		if err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Balance,
			&account.MeanPayment,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *UserRep) CheckUser(ctx context.Context, userID uuid.UUID) error {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);`, userID).Scan(&exists)

	if err != nil {
		return fmt.Errorf("[repository] failed request checkUser %w", err)
	}

	if !exists {
		return fmt.Errorf("[repo] %w: %v", &models.NoSuchUserError{UserID: userID}, err)
	}

	return nil
}

func (r *UserRep) UpdateUser(ctx context.Context, user *models.User) error { // need test
	query := `UPDATE users
				   SET username = $2,
				       planned_budget = $3,
					   avatar_url = $4
				   WHERE id = $1;`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.PlannedBudget, user.AvatarURL)
	if err != nil {
		return fmt.Errorf("[repo] failed update user %w", err)
	}

	return nil
}

func (r *UserRep) IsLoginUnique(ctx context.Context, login string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE login = $1"
	var count int
	err := r.db.QueryRowContext(ctx, query, login).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("[repo] failed login unique check %w", err)
	}

	return count == 0, nil
}

func (r *UserRep) UpdatePhoto(ctx context.Context, userID uuid.UUID, path uuid.UUID) error {
	query := `UPDATE users
				SET avatar_url = $2
			  WHERE id = $1;`
	_, err := r.db.ExecContext(ctx, query, userID, path)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("[repo] %w: %v", &models.NoSuchUserError{UserID: userID}, err)
	} else if err != nil {
		return fmt.Errorf("[repo] failed request db %s, %w", query, err)
	}
	return nil
}
