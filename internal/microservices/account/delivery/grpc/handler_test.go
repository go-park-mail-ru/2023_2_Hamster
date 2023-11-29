package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	proto "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account/delivery/grpc/generated"
	mocks "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedAccountID := uuid.New()

	request := &proto.CreateRequest{
		Id:             expectedAccountID.String(),
		UserId:         expectedAccountID.String(),
		Balance:        100.0,
		Accumulation:   true,
		BalanceEnabled: true,
		MeanPayment:    "monthly",
	}

	mockAccountServices := mocks.NewMockUsecase(ctrl)

	mockAccountServices.EXPECT().
		CreateAccount(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(expectedAccountID, nil)

	accountGRPC := NewAccountGRPC(mockAccountServices, *logger.NewLogger(context.TODO()))

	response, err := accountGRPC.Create(context.Background(), request)

	assert.NoError(t, err)

	assert.NotNil(t, response)

}

func TestUpdateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedAccountID := uuid.New()

	request := &proto.UpdasteRequest{
		Id:             expectedAccountID.String(),
		UserId:         expectedAccountID.String(),
		Balance:        150.0,
		Accumulation:   false,
		BalanceEnabled: true,
		MeanPayment:    "weekly",
	}

	mockAccountServices := mocks.NewMockUsecase(ctrl)

	mockAccountServices.EXPECT().
		UpdateAccount(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	accountGRPC := NewAccountGRPC(mockAccountServices, *logger.NewLogger(context.TODO()))

	_, err := accountGRPC.Update(context.Background(), request)

	assert.NoError(t, err)
}

func TestDeleteAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedAccountID := uuid.New()

	request := &proto.DeleteRequest{
		AccountId: expectedAccountID.String(),
		UserId:    expectedAccountID.String(),
	}

	mockAccountServices := mocks.NewMockUsecase(ctrl)

	mockAccountServices.EXPECT().
		DeleteAccount(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(errors.New("account not found"))

	accountGRPC := NewAccountGRPC(mockAccountServices, *logger.NewLogger(context.TODO()))

	_, err := accountGRPC.Delete(context.Background(), request)

	assert.Error(t, err)
}
