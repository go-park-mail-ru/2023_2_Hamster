package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	proto "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/delivery/grpc/generated"
	mocks "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/mocks"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedTagID := uuid.New()

	request := &proto.CreateTagRequest{
		UserId:      expectedTagID.String(),
		ParentId:    expectedTagID.String(),
		Name:        "TestTag",
		ShowIncome:  true,
		ShowOutcome: true,
		Regular:     false,
	}

	mockCategoryServices := mocks.NewMockUsecase(ctrl)
	mockCategoryServices.EXPECT().
		CreateTag(gomock.Any(), gomock.Any()).
		Return(expectedTagID, nil)

	categoryGRPC := NewCategoryGRPC(mockCategoryServices, *logger.NewLogger(context.TODO()))

	response, err := categoryGRPC.CreateTag(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedTagID.String(), response.TagId)
}

func TestGetTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUserID := uuid.New()

	request := &proto.UserIdRequest{
		UserId: expectedUserID.String(),
	}

	expectedTags := []models.Category{{}, {}}

	mockCategoryServices := mocks.NewMockUsecase(ctrl)
	mockCategoryServices.EXPECT().
		GetTags(gomock.Any(), gomock.Any()).
		Return(expectedTags, nil)

	categoryGRPC := NewCategoryGRPC(mockCategoryServices, *logger.NewLogger(context.TODO()))

	response, err := categoryGRPC.GetTags(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Categories, len(expectedTags))

}

func TestGetTags2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUserID := uuid.New()

	request := &proto.UserIdRequest{
		UserId: expectedUserID.String(),
	}

	expectedError := errors.New("some error message")

	mockCategoryServices := mocks.NewMockUsecase(ctrl)
	mockCategoryServices.EXPECT().
		GetTags(gomock.Any(), gomock.Any()).
		Return(nil, expectedError)

	categoryGRPC := NewCategoryGRPC(mockCategoryServices, *logger.NewLogger(context.TODO()))

	response, err := categoryGRPC.GetTags(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, expectedError, err)
}

func TestUpdateTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedID := uuid.New()
	expectedUserID := uuid.New()
	expectedParentID := uuid.New()

	request := &proto.Category{
		Id:          expectedID.String(),
		UserId:      expectedUserID.String(),
		ParentId:    expectedParentID.String(),
		Name:        "UpdatedTag",
		ShowIncome:  true,
		ShowOutcome: false,
		Regular:     true,
	}

	mockCategoryServices := mocks.NewMockUsecase(ctrl)
	mockCategoryServices.EXPECT().
		UpdateTag(gomock.Any(), gomock.Any()).
		Return(nil) // Предполагаем, что обновление прошло успешно.

	categoryGRPC := NewCategoryGRPC(mockCategoryServices, *logger.NewLogger(context.TODO()))

	response, err := categoryGRPC.UpdateTag(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Проверим, что возвращенные данные соответствуют ожидаемым.
	assert.Equal(t, expectedID.String(), response.Id)
	assert.Equal(t, expectedUserID.String(), response.UserId)
	assert.Equal(t, expectedParentID.String(), response.ParentId)
	assert.Equal(t, request.Name, response.Name)
	assert.Equal(t, request.ShowIncome, response.ShowIncome)
	assert.Equal(t, request.ShowOutcome, response.ShowOutcome)
	assert.Equal(t, request.Regular, response.Regular)
}

func TestDeleteTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedTagID := uuid.New()
	expectedUserID := uuid.New()

	request := &proto.DeleteRequest{
		TagId:  expectedTagID.String(),
		UserId: expectedUserID.String(),
	}

	mockCategoryServices := mocks.NewMockUsecase(ctrl)
	mockCategoryServices.EXPECT().
		DeleteTag(gomock.Any(), expectedTagID, expectedUserID).
		Return(nil) // Предполагаем, что удаление прошло успешно.

	categoryGRPC := NewCategoryGRPC(mockCategoryServices, *logger.NewLogger(context.TODO()))

	response, err := categoryGRPC.DeleteTag(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Проверим, что возвращенные данные соответствуют ожидаемым.
	assert.Equal(t, &empty.Empty{}, response)
}
