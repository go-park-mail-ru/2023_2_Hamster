package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/question"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Usecase struct {
	questionRep question.Repository
}

func NewUsecase(
	qr question.Repository) *Usecase {
	return &Usecase{
		questionRep: qr,
	}
}

func (uc *Usecase) CreateAnswer(ctx context.Context, userID uuid.UUID, a models.Answer) error {
	err := uc.questionRep.CreateAnswer(ctx, userID, a)
	if err != nil {
		return fmt.Errorf("error submitting answer: %w", err)
	}
	return nil
}

func (uc *Usecase) CheckUserAnswer(ctx context.Context, userID uuid.UUID, questionName string) (bool, error) {
	answerBool, err := uc.questionRep.CheckUserAnswer(ctx, userID, questionName)
	if err != nil {
		return answerBool, fmt.Errorf("error checkAnswer answer: %w", err)
	}
	return answerBool, err
}

func (uc *Usecase) CalculateAverageRating(ctx context.Context, questionName string) (float64, error) {
	answerInt, err := uc.questionRep.CalculateAverageRating(ctx, questionName)
	if err != nil {
		return answerInt, fmt.Errorf("error average rating answer: %w", err)
	}
	return answerInt, err
}
