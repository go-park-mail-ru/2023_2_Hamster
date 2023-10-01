package postgresql

import (
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
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

func (r *UserRep) CreateUser() {
	userID := uuid.New()

	_, err := r.db.ExecContext(ctx)
}

func (r *AuthRepo) CheckUser() {

}
