package sessions

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Usecase interface {
	GetSessionByCookie(ctx context.Context, cookie string) (models.Session, error)
	CreateSessionById(ctx context.Context, userID uuid.UUID) (models.Session, error)
	DeleteSessionByCookie(ctx context.Context, cookie string) error
}

type Repository interface {
	GetSessionByCookie(ctx context.Context, cookie string) (models.Session, error)
	CreateSession(ctx context.Context, session models.Session) error
	DeleteSession(ctx context.Context, cookie string) error
}
