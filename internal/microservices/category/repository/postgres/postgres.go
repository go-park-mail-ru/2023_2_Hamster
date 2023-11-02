package postgres

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgtype/pgxtype"
)

const (
	CreateCategory       = `INSERT INTO users (login, username, password_hash) VALUES ($1, $2, $3) RETURNING id;`
	UserIDGetByID        = `SELECT * FROM users WHERE id = $1;`
	UserGetByUserName    = `SELECT id, login, username, password_hash, planned_budget, avatar_url From users WHERE (login=$1)`
	UserGetPlannedBudget = "SELECT planned_budget FROM users WHERE id = $1"
	UserCheck            = `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);`
	UserUpdate           = `UPDATE users SET username = $2, planned_budget = $3, avatar_url = $4 WHERE id = $1;`
	UserUpdatePhoto      = `UPDATE users SET avatar_url = $2 WHERE id = $1;`
	AccountBalance       = "SELECT SUM(balance) FROM accounts WHERE user_id = $1" // TODO: move accounts
	AccountGet           = `SELECT * FROM accounts WHERE user_id = $1`            // TODO: move accounts
)

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

func (r *Repository) CreateCategory(ctx context.Context, category models.Category) (uuid.UUID, error) {

}

func (r *Repository) UpdateCategory(ctx context.Context, transaction *models.Category) error {

}

func (r *Repository) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {

}

func (r *Repository) GetFeed(ctx context.Context, userID uuid.UUID) ([]models.Category, error) {

}
