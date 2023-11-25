package grpc

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"

	proto "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/grpc/generated"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authGRPC struct {
	authServices auth.Usecase
	// sessionService sessions.Usecase
	logger logger.Logger

	proto.UnimplementedAuthServiceServer
}

func NewAuthGRPC(authServices auth.Usecase, log logger.Logger) *authGRPC {
	return &authGRPC{
		authServices: authServices,
		logger:       log,
	}
}

func (a *authGRPC) SignUp(ctx context.Context, in *proto.SignUpRequest) (*proto.SignUpResponse, error) {
	request := auth.SignUpInput{
		Login:          in.Login,
		Username:       in.Username,
		PlaintPassword: in.Password,
	}

	id, login, username, err := a.authServices.SignUp(ctx, request)
	if err != nil {
		var errUserAlreadyExists *models.UserAlreadyExistsError
		if errors.As(err, &errUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	body := proto.SignUpResponseBody{
		Id:       id.String(),
		Login:    login,
		Username: username,
	}

	return &proto.SignUpResponse{
		Status: "200",
		Body:   &body,
	}, nil
}

func (a *authGRPC) Login(ctx context.Context, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	id, login, username, err := a.authServices.Login(ctx, in.Login, in.Password)
	if err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		var errIncorrectPassword *models.IncorrectPasswordError
		if errors.As(err, &errIncorrectPassword) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	body := proto.LoginResponseBody{
		Id:       id.String(),
		Login:    login,
		Username: username,
	}

	return &proto.LoginResponse{
		Status: "200",
		Body:   &body,
	}, nil
}

func (a *authGRPC) CheckLoginUnique(ctx context.Context, in *proto.UniqCheckRequest) (*proto.UniqCheckResponse, error) {
	isUniq, err := a.authServices.CheckLoginUnique(ctx, in.Login)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.UniqCheckResponse{
		Status: "200",
		Body:   !isUniq,
	}, nil
}

func (a *authGRPC) GetByID(ctx context.Context, in *proto.UserIdRequest) (*proto.UserResponse, error) {
	userUUID, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	user, err := a.authServices.GetByID(ctx, userUUID)
	if err != nil {
		var errNoSuchUser models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.UserResponse{
		Status: "200",
		Body: &proto.UserResponseBody{
			Id:       user.ID.String(),
			Login:    user.Login,
			Username: user.Username,
		},
	}, nil
}

// func (a *authGRPC) HealthCheck(ctx context.Context, in *emptypb.Empty) (*proto.HelthCheckResponse, error) {
// 	/// a.sessionService.GetSessionByCookie(ctx)
// 	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
// }
// func (a *authGRPC) LogOut(ctx context.Context, in *emptypb.Empty) (*proto.LogoutResponse, error) {
// 	return nil, status.Errorf(codes.Unimplemented, "method LogOut not implemented")
// }
