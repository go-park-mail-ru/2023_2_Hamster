package postgresql

import (
	"database/sql"

	"github.com/google/uuid"
)

const (
	userExist = "smth"
)

type AuthRepo struct {
	db *sql.DB
}

func userRepository(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		db:     db,
		logger: l,
	}
}

func (r *AuthRepo) CreateUser() {
	userID := uuid.New()

	_, err := r.db.ExecContext(ctx)
}

func (r *AuthRepo) CheckUser() {

}
