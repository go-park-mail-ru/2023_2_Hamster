package postgresql

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/jackc/pgx/v5"
)

const (
	UserCheckLoginUnique = "SELECT COUNT(*) FROM users WHERE login = $1"
)

type AuthRep struct {
	db     *pgx.Conn
	logger logger.CustomLogger
}

func NewRepository(db *pgx.Conn, l logger.CustomLogger) *AuthRep {
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
