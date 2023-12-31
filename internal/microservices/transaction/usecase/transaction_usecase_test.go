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
		name                string
		expectedTransaction []models.Transaction
		expectedErr         error
		mockRepoFn          func(*mock.MockRepository)
	}{
		{
			name:                "Successful TestUsecase_GetFeed",
			expectedTransaction: []models.Transaction{},
			expectedErr:         nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().GetFeed(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Transaction{}, nil)
			},
		},
		{
			name:                "Error in TestUsecase_GetFeed",
			expectedTransaction: []models.Transaction{},
			expectedErr:         fmt.Errorf("[usecase] can't get transactions from repository some error"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().GetFeed(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Transaction{}, errors.New("some error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			tc.mockRepoFn(mockRepo)

			mockUsecase := NewUsecase(mockRepo, *logger.NewLogger(context.TODO()))

			userID := uuid.New()

			transaciton, err := mockUsecase.GetFeed(context.Background(), userID, &models.QueryListOptions{})

			assert.Equal(t, tc.expectedTransaction, transaciton)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_GetCount(t *testing.T) {
	testCases := []struct {
		name                string
		expectedTransaction int
		expectedErr         error
		mockRepoFn          func(*mock.MockRepository)
	}{
		{
			name:                "Successful TestUsecase_GetCount",
			expectedTransaction: 1,
			expectedErr:         nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().GetCount(gomock.Any(), gomock.Any()).Return(1, nil)
			},
		},
		{
			name:                "Error in TestUsecase_GetCount",
			expectedTransaction: 1,
			expectedErr:         fmt.Errorf("[usecase] can't get count transactions from repository some error"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().GetCount(gomock.Any(), gomock.Any()).Return(1, errors.New("some error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			tc.mockRepoFn(mockRepo)

			mockUsecase := NewUsecase(mockRepo, *logger.NewLogger(context.TODO()))

			userID := uuid.New()

			transaciton, err := mockUsecase.GetCount(context.Background(), userID)

			assert.Equal(t, tc.expectedTransaction, transaciton)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_CreateTransaction(t *testing.T) {
	userIdTest := uuid.New()
	testCases := []struct {
		name                  string
		expectedTransactionID uuid.UUID
		expectedErr           error
		mockRepoFn            func(*mock.MockRepository)
	}{
		{
			name:                  "Successful TestUsecase_CreateTransaction",
			expectedTransactionID: userIdTest,
			expectedErr:           nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(userIdTest, nil)
			},
		},
		{
			name:                  "Error in TestUsecase_CreateTransaction",
			expectedErr:           fmt.Errorf("[usecase] can't create transaction into repository: some error"),
			expectedTransactionID: userIdTest,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(userIdTest, errors.New("some error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			tc.mockRepoFn(mockRepo)

			mockUsecase := NewUsecase(mockRepo, *logger.NewLogger(context.TODO()))

			transaction := models.Transaction{}
			transactionID, err := mockUsecase.CreateTransaction(context.Background(), &transaction)
			assert.Equal(t, tc.expectedTransactionID, transactionID)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_UpdateTransaction(t *testing.T) {
	userIdTest := uuid.New()
	testCases := []struct {
		name        string
		expectedErr error
		mockRepoFn  func(*mock.MockRepository)
	}{
		{
			name:        "Successful",
			expectedErr: nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckForbidden(gomock.Any(), gomock.Any()).Return(userIdTest, nil)
				mockRepositry.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:        "Error in userIDCheck != transaction.UserID",
			expectedErr: fmt.Errorf("[usecase] can't be update by user: user has no rights"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckForbidden(gomock.Any(), gomock.Any()).Return(uuid.New(), nil)
				//mockRepositry.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(userIdTest, errors.New("some error"))
			},
		},
		{
			name:        "Error in can't find transaction in repository",
			expectedErr: fmt.Errorf("[usecase] can't find transaction in repository some err"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckForbidden(gomock.Any(), gomock.Any()).Return(uuid.New(), errors.New("some err"))
				//mockRepositry.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(userIdTest, errors.New("some error"))
			},
		},
		{
			name:        "Error in can't find transaction in repository",
			expectedErr: fmt.Errorf("[usecase] can't update transaction some error"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckForbidden(gomock.Any(), gomock.Any()).Return(userIdTest, nil)
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

			mockUsecase := NewUsecase(mockRepo, *logger.NewLogger(context.TODO()))

			transaction := models.Transaction{UserID: userIdTest}
			err := mockUsecase.UpdateTransaction(context.Background(), &transaction)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_DeleteTransaction(t *testing.T) {
	userIdTest := uuid.New()
	testCases := []struct {
		name        string
		expectedErr error
		mockRepoFn  func(*mock.MockRepository)
	}{
		{
			name:        "Successful",
			expectedErr: nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckForbidden(gomock.Any(), gomock.Any()).Return(userIdTest, nil)
				mockRepositry.EXPECT().DeleteTransaction(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:        "Error in userIDCheck != transaction.UserID",
			expectedErr: fmt.Errorf("[usecase] can't be deleted by user: user has no rights"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckForbidden(gomock.Any(), gomock.Any()).Return(uuid.New(), nil)
				//mockRepositry.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(userIdTest, errors.New("some error"))
			},
		},
		{
			name:        "Error in can't find transaction in repository",
			expectedErr: fmt.Errorf("[usecase] can't find transaction in repository some err"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckForbidden(gomock.Any(), gomock.Any()).Return(uuid.New(), errors.New("some err"))
				//mockRepositry.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(userIdTest, errors.New("some error"))
			},
		},
		{
			name:        "Error in can't find transaction in repository",
			expectedErr: fmt.Errorf("[usecase] can`t be deleted from repository"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckForbidden(gomock.Any(), gomock.Any()).Return(userIdTest, nil)
				mockRepositry.EXPECT().DeleteTransaction(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("some error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			tc.mockRepoFn(mockRepo)

			mockUsecase := NewUsecase(mockRepo, *logger.NewLogger(context.TODO()))

			err := mockUsecase.DeleteTransaction(context.Background(), userIdTest, userIdTest)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_GetTransactionForExport(t *testing.T) {
	testCases := []struct {
		name           string
		expectedResult []models.TransactionExport
		expectedErr    error
		mockRepoFn     func(*mock.MockRepository)
	}{
		{
			name:           "Successful GetTransactionForExport",
			expectedResult: []models.TransactionExport{},
			expectedErr:    nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().GetTransactionForExport(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.TransactionExport{}, nil)
			},
		},
		{
			name:           "Error in GetTransactionForExport",
			expectedResult: nil,
			expectedErr:    fmt.Errorf("[usecase] can't get transactions from repository some error"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().GetTransactionForExport(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			tc.mockRepoFn(mockRepo)

			mockUsecase := NewUsecase(mockRepo, *logger.NewLogger(context.TODO()))

			userID := uuid.New()
			query := &models.QueryListOptions{}

			result, err := mockUsecase.GetTransactionForExport(context.Background(), userID, query)

			assert.Equal(t, tc.expectedResult, result)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}
