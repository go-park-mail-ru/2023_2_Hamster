package grpc

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/question"
	proto "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/question/delivery/grpc/generated"
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
	request := question.
	return nil, status.Errorf(codes.Unimplemented, "method CreateAnswer not implemented")
}
func (q *questionGRPC) CheckUserAnswer(context.Context, *proto.CheckUserAnswerResponse) (*proto.CheckUserAnswerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckUserAnswer not implemented")
}
func (q *questionGRPC) CalculateAverageRating(context.Context, *proto.QuestionNameRequest) (*proto.AverageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CalculateAverageRating not implemented")
}
