package postgresql

import (

	"database/sql"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
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

func (r *AuthRepo) CheckUser(username string) (bool, error) {
	var exists bool
	query := `SELECT exists(SELECT 1 FROM users WHERE username=\$1)`
	err := r.db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil

