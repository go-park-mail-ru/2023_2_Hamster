package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	mock "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/transaction/mocks"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_GetFeed(t *testing.T) {
	testCases := []struct {
		name                      string
		expectedTransactionModels []models.Transaction
		expectedErr               error
		mockRepoFn                func(*mock.MockRepository)
	}{
		{
			name:                      "Successful get feed",
			expectedTransactionModels: []models.Transaction{},
			expectedErr:               nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().GetFeed(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Transaction{}, true, nil)
			},
		},
		{
			name:                      "Error in get feed",
			expectedTransactionModels: []models.Transaction{},
			expectedErr:               fmt.Errorf("[usecase] can't get transactions from repository some error"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().GetFeed(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Transaction{}, true, errors.New("some error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			tc.mockRepoFn(mockRepo)

			mockUsecase := NewUsecase(mockRepo, *logger.CreateCustomLogger())

			userID := uuid.New()

			transaction, _, err := mockUsecase.GetFeed(context.Background(), userID, 10, 5)

			assert.Equal(t, tc.expectedTransactionModels, transaction)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_CreateTransaction(t *testing.T) {
	testUUID := uuid.New()
	testCases := []struct {
		name                      string
		expectedTransactionModels uuid.UUID
		expectedErr               error
		mockRepoFn                func(*mock.MockRepository)
	}{
		{
			name:                      "Successful create transaciton",
			expectedTransactionModels: testUUID,
			expectedErr:               nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(testUUID, nil)
			},
		},
		{
			name:                      "Error in create transaction",
			expectedTransactionModels: testUUID,
			expectedErr:               fmt.Errorf("[usecase] can't create transaction into repository: some error"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(testUUID, errors.New("some error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			tc.mockRepoFn(mockRepo)

			mockUsecase := NewUsecase(mockRepo, *logger.CreateCustomLogger())

			transactionM := models.Transaction{}
			transaction, err := mockUsecase.CreateTransaction(context.Background(), &transactionM)

			assert.Equal(t, tc.expectedTransactionModels, transaction)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_UpdateTransaction(t *testing.T) {
	testUUID := uuid.New()
	testCases := []struct {
		name                      string
		expectedTransactionModels uuid.UUID
		expectedErr               error
		mockRepoFn                func(*mock.MockRepository)
	}{
		{
			name:                      "Successful update transaciton",
			expectedTransactionModels: testUUID,
			expectedErr:               fmt.Errorf("[usecase] can't create transaction into repository: some error"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckForbidden(gomock.Any(), gomock.Any()).Return(testUUID, nil)

				mockRepositry.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Return(errors.New("fasdf"))

			},
		},
		{
			name:                      "Error in update transaction",
			expectedTransactionModels: testUUID,
			expectedErr:               fmt.Errorf("[usecase] can't create transaction into repository: some error"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckForbidden(gomock.Any(), gomock.Any()).Return(testUUID, nil)
				mockRepositry.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Return(errors.New("some error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			tc.mockRepoFn(mockRepo)

			mockUsecase := NewUsecase(mockRepo, *logger.CreateCustomLogger())

			transactionM := models.Transaction{}
			err := mockUsecase.UpdateTransaction(context.Background(), &transactionM)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

// func TestUsecase_DeleteTransaction(t *testing.T) {
// 	testCases := []struct {
// 		name                      string
// 		expectedTransactionModels []models.Transaction
// 		expectedErr               error
// 		mockRepoFn                func(*mock.MockRepository)
// 	}{
// 		{
// 			name:                      "Successful get feed",
// 			expectedTransactionModels: []models.Transaction{},
// 			expectedErr:               fmt.Errorf("[usecase] can't be update by user: user has no rights"),
// 			mockRepoFn: func(mockRepositry *mock.MockRepository) {
// 				mockRepositry.EXPECT().GetFeed(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Transaction{}, true, nil)
// 			},
// 		},
// 		{
// 			name:                      "Error in get feed",
// 			expectedTransactionModels: []models.Transaction{},
// 			expectedErr:               fmt.Errorf("[usecase] can't be update by user: user has no rights"),
// 			mockRepoFn: func(mockRepositry *mock.MockRepository) {
// 				mockRepositry.EXPECT().GetFeed(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Transaction{}, true, errors.New("some error"))
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockRepo := mock.NewMockRepository(ctrl)
// 			tc.mockRepoFn(mockRepo)

// 			mockUsecase := NewUsecase(mockRepo, *logger.CreateCustomLogger())

// 			userID := uuid.New()

// 			transaction, _, err := mockUsecase.GetFeed(context.Background(), userID, 10, 5)

// 			assert.Equal(t, tc.expectedTransactionModels, transaction)
// 			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
// 				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
// 			}
// 		})
// 	}
// }
