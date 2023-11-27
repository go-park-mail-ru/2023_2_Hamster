package grpc

import (
	"context"

	proto "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account"
)

type accountGRPC struct {
	AccountServices account.Usecase
	logger          logger.Logger

	proto.UnimplementedAccountServiceServer
}

func NewAccountGRPC(accountServices account.Usecase, logger logger.Logger) *accountGRPC {
	return &accountGRPC{
		AccountServices: accountServices,
		logger:          logger,
	}
}

func (a *accountGRPC) Create(ctx context.Context, in *proto.CreateRequest) (*proto.CreateAccountResponse, error) {
	uuidID, _ := uuid.Parse(in.Id)
	request := models.Accounts{
		ID:             uuidID,
		Balance:        float64(in.Balance),
		Accumulation:   in.Accumulation,
		BalanceEnabled: in.BalanceEnabled,
		MeanPayment:    in.MeanPayment,
	}
	userID, _ := uuid.Parse(in.UserId)
	accountID, err := a.AccountServices.CreateAccount(ctx, userID, &request)

	return &proto.CreateAccountResponse{AccountId: accountID.String()}, err
}

func (a *accountGRPC) Update(ctx context.Context, in *proto.UpdasteRequest) (*empty.Empty, error) {
	uuidID, _ := uuid.Parse(in.Id)
	request := models.Accounts{
		ID:             uuidID,
		Balance:        float64(in.Balance),
		Accumulation:   in.Accumulation,
		BalanceEnabled: in.BalanceEnabled,
		MeanPayment:    in.MeanPayment,
	}
	userID, _ := uuid.Parse(in.UserId)

	err := a.AccountServices.UpdateAccount(ctx, userID, &request)

	return &empty.Empty{}, err
}

func (a *accountGRPC) Delete(ctx context.Context, in *proto.DeleteRequest) (*empty.Empty, error) {
	AccountUUID, _ := uuid.Parse(in.AccountId)
	userID, _ := uuid.Parse(in.UserId)

	err := a.AccountServices.DeleteAccount(ctx, userID, AccountUUID)

	return &empty.Empty{}, err
}
