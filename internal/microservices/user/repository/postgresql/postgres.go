package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

const (
	SharingUserGet = `SELECT
						u.id,
						u.login,
						u.avatar_url
					FROM Users u
					JOIN UserAccount ua ON u.id = ua.user_id
					WHERE ua.account_id = $1;`

	UserCreate           = `INSERT INTO users (login, username, password_hash) VALUES ($1, $2, $3) RETURNING id;`
	UserIDGetByID        = `SELECT id, login, username, password_hash, planned_budget, avatar_url FROM users WHERE id = $1;`
	UserGetByUserName    = `SELECT id, login, username, password_hash, planned_budget, avatar_url From users WHERE (login=$1)`
	UserGetPlannedBudget = "SELECT planned_budget FROM users WHERE id = $1"
	UserCheck            = `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);`
	UserUpdate           = `UPDATE users SET username = $2, planned_budget = $3, avatar_url = $4 WHERE id = $1;`
	UserUpdatePhoto      = `UPDATE users SET avatar_url = $2 WHERE id = $1;`
	AccountBalance       = `SELECT SUM(a.balance)
							FROM Accounts a
							JOIN UserAccount ua ON a.id = ua.account_id
							WHERE ua.user_id = $1 
							AND a.balance_enabled = true;` // TODO: move accounts

	AccountGet = `SELECT a.*
				  FROM Accounts a
				  JOIN UserAccount ua ON a.id = ua.account_id
				  WHERE ua.user_id = $1` // TODO: move accounts

	ActualBudgetCalculation = `SELECT SUM(outcome) AS total_sum
								FROM transaction
								WHERE date_part('month', date) = date_part('month', CURRENT_DATE)
								AND date_part('year', date) = date_part('year', CURRENT_DATE)
								AND outcome > 0
								AND account_income = account_outcome
								AND user_id = $1;`
)

type UserRep struct {
	db     postgresql.DbConn
	logger logger.Logger
}

func NewRepository(db postgresql.DbConn, log logger.Logger) *UserRep {
	return &UserRep{
		db:     db,
		logger: log,
	}
}

func (r *UserRep) CreateUser(ctx context.Context, u models.User) (uuid.UUID, error) {
	row := r.db.QueryRow(ctx, UserCreate, u.Login, u.Username, u.Password)
	var id uuid.UUID

	err := row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("error request %w", err)
	}
	return id, nil
}

func (r *UserRep) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	row := r.db.QueryRow(ctx, UserIDGetByID, userID)
	var u models.User

	err := row.Scan(&u.ID, &u.Login, &u.Username, &u.Password, &u.PlannedBudget, &u.AvatarURL)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("[repo] %w: %v", &models.NoSuchUserError{UserID: userID}, err)
	} else if err != nil {
		return nil,
			fmt.Errorf("failed request db %s, %w", UserIDGetByID, err)

	}
	return nil, nil
}

func (r *UserRep) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	row := r.db.QueryRow(ctx, UserGetByUserName, login)
	var u models.User
	err := row.Scan(&u.ID, &u.Login, &u.Username, &u.Password, &u.PlannedBudget, &u.AvatarURL)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("[repo] nothing found for this request %w: %v", &models.NoSuchUserInLogin{Login: login}, err)
	} else if err != nil {
		return nil,
			fmt.Errorf("[repo] failed request db %w", err)
	}
	return &u, nil
}

func (r *UserRep) GetUserBalance(ctx context.Context, userID uuid.UUID) (float64, error) { // need check
	var totalBalance sql.NullFloat64
	err := r.db.QueryRow(ctx, AccountBalance, userID).Scan(&totalBalance)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("[repo] %w: %v", &models.NoSuchUserIdBalanceError{UserID: userID}, err)
	} else if err != nil {
		return 0, fmt.Errorf("[repo] failed request db %w", err)
	}

	if totalBalance.Valid {
		return totalBalance.Float64, nil
	}

	return 0, nil
}

func (r *UserRep) GetPlannedBudget(ctx context.Context, userID uuid.UUID) (float64, error) { // need check
	var plannedBudget sql.NullFloat64
	err := r.db.QueryRow(ctx, UserGetPlannedBudget, userID).Scan(&plannedBudget)

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

func (r *UserRep) GetCurrentBudget(ctx context.Context, userID uuid.UUID) (float64, error) { // need check
	var currentBudget sql.NullFloat64

	err := r.db.QueryRow(ctx, ActualBudgetCalculation, userID).Scan(&currentBudget)

	if err != nil {
		return 0, fmt.Errorf("[repository] failed request db %w", err)
	}

	if currentBudget.Valid {
		return currentBudget.Float64, nil
	}
	return 0, nil
}

func (r *UserRep) GetAccounts(ctx context.Context, user_id uuid.UUID) ([]models.Accounts, error) {
	var accounts []models.Accounts

	rows, err := r.db.Query(ctx, AccountGet, user_id)
	if err != nil {
		return nil, fmt.Errorf("[repo] %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var account models.Accounts

		if err := rows.Scan(
			&account.ID,
			&account.Balance,
			&account.SharingID,
			&account.Accumulation,
			&account.BalanceEnabled,
			&account.MeanPayment,
		); err != nil {
			return nil, fmt.Errorf("[repo] %w", err)
		}

		sharingUsers, err := r.getSharingUsers(ctx, account.ID)
		if err != nil {
			return nil, fmt.Errorf("[repo] %w", err)
		}

		account.Users = sharingUsers
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[repo] %w", err)
	}

	if len(accounts) == 0 {
		return nil, fmt.Errorf("[repo] %w: %v", &models.NoSuchAccounts{UserID: user_id}, err)
	}

	return accounts, nil
}

func (r *UserRep) getSharingUsers(ctx context.Context, accountID uuid.UUID) ([]models.SharingUser, error) {
	var sharingUsers []models.SharingUser

	rows, err := r.db.Query(ctx, SharingUserGet, accountID)
	if err != nil {
		return nil, fmt.Errorf("[repo] %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sharingUser models.SharingUser

		if err := rows.Scan(
			&sharingUser.ID,
			&sharingUser.Login,
			&sharingUser.AvatarURL,
		); err != nil {
			return nil, fmt.Errorf("[repo] %w", err)
		}

		sharingUsers = append(sharingUsers, sharingUser)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[repo] %w", err)
	}

	return sharingUsers, nil
}

// func (r *UserRep) CheckUser(ctx context.Context, userID uuid.UUID) error {
// 	var exists bool
// 	err := r.db.QueryRow(ctx, UserCheck, userID).Scan(&exists)

// 	if err != nil {
// 		return fmt.Errorf("[repo] failed request checkUser %w", err)
// 	}

// 	if !exists {
// 		return fmt.Errorf("[repo] %w: %v", &models.NoSuchUserError{UserID: userID}, err)
// 	}

// 	return nil
// }

func (r *UserRep) UpdateUser(ctx context.Context, user *models.User) error { // need test
	_, err := r.db.Exec(ctx, UserUpdate, user.ID, user.Username, user.PlannedBudget, user.AvatarURL)
	if err != nil {
		return fmt.Errorf("[repo] failed update user %w", err)
	}

	return nil
}

func (r *UserRep) UpdatePhoto(ctx context.Context, userID uuid.UUID, path uuid.UUID) error {
	_, err := r.db.Exec(ctx, UserUpdatePhoto, userID, path)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("[repo] %w: %v", &models.NoSuchUserError{UserID: userID}, err)
	} else if err != nil {
		return fmt.Errorf("[repo] failed request db: %s, %w", UserUpdatePhoto, err)
	}
	return nil
}
