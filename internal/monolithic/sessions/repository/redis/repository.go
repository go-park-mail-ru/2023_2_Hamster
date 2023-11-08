package redis

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionRep struct {
	db *redis.Client
}

func NewSessionRepository(db *redis.Client) *SessionRep {
	return &SessionRep{
		db: db,
	}
}

func (r *SessionRep) GetSessionByCookie(ctx context.Context, cookie string) (models.Session, error) {
	var session models.Session
	result, err := r.db.Get(context.TODO(), cookie).Result()
	if err == redis.Nil {
		return models.Session{}, sessions.ErrSessionNotFound
	}

	session.UserId, err = uuid.Parse(result)
	if err != nil {
		return models.Session{}, sessions.ErrInvalidUUID
	}

	session.Cookie = cookie

	return session, nil
}

func (r *SessionRep) CreateSession(ctx context.Context, session models.Session) error {
	err := r.db.Set(context.TODO(), session.Cookie, session.UserId.String(), 0).Err()
	return err
}

func (r *SessionRep) DeleteSession(ctx context.Context, cookie string) error {
	err := r.db.Del(context.TODO(), cookie).Err()
	if err == redis.Nil {
		return sessions.ErrSessionNotFound
	}

	return err
}
