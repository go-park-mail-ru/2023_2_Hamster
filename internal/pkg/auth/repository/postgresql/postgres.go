package postgresql

import (
	"database/sql"
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
}

func (r *AuthRepo) CheckUser() {
	panic("unimplemented")
}
