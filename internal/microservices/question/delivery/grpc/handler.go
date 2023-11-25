package grpc

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/question"
	proto "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/question/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type questionGRPC struct {
	questoinServices question.Usecase
	logger           logger.Logger

	proto.UnimplementedQuestionServiceServer
}

func NewAuthGRPC(questionServices question.Usecase, log logger.Logger) *questionGRPC {
	return &questionGRPC{
		questoinServices: questionServices,
		logger:           log,
	}
}

func (q *questionGRPC) CreateAnswer(ctx context.Context, in *proto.AnswerRequest) (*emptypb.Empty, error) {
	request := models.Answer{
		Name:   in.Name,
		Rating: int(in.Rating),
	}
	uuidUser, _ := uuid.Parse(in.Id)
	err := q.questoinServices.CreateAnswer(ctx, uuidUser, request)
	if err != nil {
		q.logger.Error("failed in uc.CreateAnswer", err)

		return nil, status.Errorf(codes.Internal, "method CreateAnswer not implemented")
	}

	return nil, nil
}
func (q *questionGRPC) CheckUserAnswer(ctx context.Context, in *proto.CheckUserAnswerRequest) (*proto.CheckUserAnswerResponse, error) {
	request := proto.CheckUserAnswerRequest{
		Id:           in.Id,
		QuestionName: in.QuestionName,
	}
	uuidUser, _ := uuid.Parse(in.Id)
	qBool, err := q.questoinServices.CheckUserAnswer(ctx, uuidUser, request.QuestionName)
	if err != nil {
		q.logger.Error("failed in uc.CreateAnswer", err)

		return nil, status.Errorf(codes.Internal, "method CreateAnswer not implemented")

	}

	return &proto.CheckUserAnswerResponse{Average: qBool}, nil
}

func (q *questionGRPC) CalculateAverageRating(ctx context.Context, in *proto.CalculateAverageRatingRequest) (*proto.AverageResponse, error) {
	request := proto.CalculateAverageRatingRequest{
		QuestionName: in.QuestionName,
	}
	intQuestion, err := q.questoinServices.CalculateAverageRating(ctx, request.QuestionName)

	if err != nil {
		q.logger.Error("failed in uc.CreateAnswer", err)

		return nil, status.Errorf(codes.Internal, "method CreateAnswer not implemented")
	}

	return &proto.AverageResponse{Average: int32(intQuestion)}, nil
}
