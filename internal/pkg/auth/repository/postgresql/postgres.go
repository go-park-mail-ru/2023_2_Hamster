package postgresql

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	userExist = "smth"
)

type AuthRep struct {
	db     *sqlx.DB
	logger logger.CustomLogger
}

func NewRepository(db *sqlx.DB, l logger.CustomLogger) *AuthRep {
	return &AuthRep{
		db:     db,
		logger: l,
	}
}

/*func (r *AuthRepo) CreateUser() {
	userID := uuid.New()

	_, err := r.db.ExecContext(ctx)
}*/

func (r *AuthRep) CheckUser(username string) (bool, error) {
	var exists bool
	query := `SELECT exists(SELECT 1 FROM users WHERE username=\$1)`
	err := r.db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("[repository] user %s don't exist: %v", username, err)
	}
	return exists, nil
}

func (r *AuthRep) GetUserByAuthData(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	query := fmt.Sprintf(
		`SELECT 
		CASE 
			WHEN is_income = true THEN SUM(total)
			WHEN is_income = false THEN -SUM(total)
		END AS calculated_total
	FROM 
		Transaction
	WHERE 
		user_id = \$1
		AND date >= DATE_TRUNC('month', CURRENT_DATE)
		AND date < DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month';`)

	var user models.User
	err := r.db.QueryRow(query, userID).Scan(&user.ID, &user.Username,
		&user.FirstName, &user.LastName, &user.PlannedBudget, &user.Password,
		&user.AvatarURL, &user.Salt)
	if err != nil {
		return nil, fmt.Errorf("[repository] Error: %v", err)
	}
	return &user, nil
}
