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

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (r *AuthRepo) CreateUser() {
	userID := uuid.New()

	_, err := r.db.ExecContext(ctx)
}

func (r *AuthRepo) CheckUser() {
	panic("unimplemented")
}
