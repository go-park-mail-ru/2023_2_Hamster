package question

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Usecase interface {
	CreateAnswer(ctx context.Context, userID uuid.UUID, a models.Answer) error
	CheckUserAnswer(ctx context.Context, userID, questionName string) (bool, error)
	CalculateAverageRating(ctx context.Context, questionName string) (int, error)
}

type Repository interface {
	CreateAnswer(ctx context.Context, userID uuid.UUID, a models.Answer) error
	CheckUserAnswer(ctx context.Context, userID, questionName string) (bool, error)
	CalculateAverageRating(ctx context.Context, questionName string) (int, error)
}
