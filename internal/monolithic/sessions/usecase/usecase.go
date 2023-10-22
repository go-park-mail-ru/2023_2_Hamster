package usecase

import "github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions"

type Usecase struct {
	sessionRepo sessions.Repository
}

func NewSessionUsecase(sessionRepo sessions.Repository) *sessions.Usecase {
	return &Usecase{sessionRepo: sessionRepo}
}
