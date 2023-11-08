package sessions

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

var (
	ErrSessionNotFound         = errors.New("session not found")
	ErrInvalidUUID             = errors.New("uuid parse error")
	ErrSessionIsAlreadyCreated = errors.New("session already created")
	ErrInternalServer          = errors.New("internal server error")
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
