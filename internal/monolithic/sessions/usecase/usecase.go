package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/hasher"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions"
	"github.com/google/uuid"
)

type Usecase struct {
	sessionRepo sessions.Repository
}

func NewSessionUsecase(sessionRepo sessions.Repository) *Usecase {
	return &Usecase{
		sessionRepo: sessionRepo,
	}
}
func (u *Usecase) GetSessionByCookie(ctx context.Context, cookie string) (models.Session, error) {
	session, err := u.sessionRepo.GetSessionByCookie(ctx, cookie)
	return session, err
}

func (u *Usecase) CreateSessionById(ctx context.Context, userID uuid.UUID) (models.Session, error) {
	session := models.Session{
		UserId: userID,
		Cookie: hasher.GenerateSession(uuid.New().String()),
	}
	err := u.sessionRepo.CreateSession(ctx, session)
	return session, err
}

func (u *Usecase) DeleteSessionByCookie(ctx context.Context, cookie string) error {
	err := u.sessionRepo.DeleteSession(ctx, cookie)
	return err
}
