package postgresql

import (
	"database/sql"
)

const (
	userExist = "smth"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
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
}
