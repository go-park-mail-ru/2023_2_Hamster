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
)

const (
	UserCheckLoginUnique = `SELECT COUNT(*) FROM users WHERE login = $1;`
	UserGetByUserName    = `SELECT id, login, username, password_hash, planned_budget, avatar_url From users WHERE (login=$1);`
	UserCreate           = `INSERT INTO users (login, username, password_hash) VALUES ($1, $2, $3) RETURNING id;`
	UserIDGetByID        = `SELECT id, login, username, password_hash, planned_budget, avatar_url FROM users WHERE id = $1;`
)

type AuthRep struct {
	db     postgresql.DbConn
	logger logger.Logger
}

func NewRepository(db postgresql.DbConn, l logger.Logger) *AuthRep {
	return &AuthRep{
		db:     db,
		logger: l,
	}
}

func (r *AuthRep) CheckLoginUnique(ctx context.Context, login string) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, UserCheckLoginUnique, login).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("[repo] failed login unique check %w", err)
	}

	return count == 0, nil
}

func (r *AuthRep) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	row := r.db.QueryRow(ctx, UserGetByUserName, login)
	var u models.User
	err := row.Scan(&u.ID, &u.Login, &u.Username, &u.Password, &u.PlannedBudget, &u.AvatarURL)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("[repo] nothing found for this request %w", err)
	} else if err != nil {
		return nil,
			fmt.Errorf("[repo] failed request db %w", err)
	}
	return &u, nil
}

func (r *AuthRep) CreateUser(ctx context.Context, u models.User) (uuid.UUID, error) {
	row := r.db.QueryRow(ctx, UserCreate, u.Login, u.Username, u.Password)
	var id uuid.UUID

	err := row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("error request %w", err)
	}
	return id, nil
}

func (r *AuthRep) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	row := r.db.QueryRow(ctx, UserIDGetByID, userID)
	var u models.User

	err := row.Scan(&u.ID, &u.Login, &u.Username, &u.Password, &u.PlannedBudget, &u.AvatarURL)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("[repo] %w: %v", &models.NoSuchUserError{UserID: userID}, err)
	} else if err != nil {
		return nil,
			fmt.Errorf("failed request db %s, %w", UserIDGetByID, err)

	}
	return &u, nil
}
