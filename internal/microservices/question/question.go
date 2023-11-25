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

type (
	/*
		message AnswerRequest {
		  string name = 1;
		  int32  rating = 2;
		};

		message AverageRatingResponse {
		  int32 averageRating = 1;
		};

		message CheckUserAnswerResponse {
		    bool average = 1;
		};

		message QuestionNameRequest {
		  string questionName = 1;
		};

		message AverageResponse {
		    int32 average = 1;
		};
	*/

	AnswerRequest struct {
		Name   string  `json:"name"`
		Rating float64 `json:"rating"`
	}

	AverageRatingResponse struct {
		AverageRating int `json:"avg_rating"`
	}

	QuestionNameRequest struct {
		QuestionName string `json:"name"`
	}

	CheckUserAnswerResponse struct {
		Average bool `json:"avg"`
	}

	AverageResponse struct {
		Average float64 `json:"avg"`
	}
)
