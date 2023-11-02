package postgres

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgtype/pgxtype"
)

const ()

type Repository struct {
	db  pgxtype.Querier
	log logger.Logger
}

func NewRepository(db pgxtype.Querier, log logger.Logger) *Repository {
	return &Repository{
		db:  db,
		log: log,
	}
}

func (rep *Repository) CreateCategory(ctx context.Context, category models.Category) (uuid.UUID, error) {

}
