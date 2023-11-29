package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	mock "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account/mocks"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_CreateAccount(t *testing.T) {
	testUUID := uuid.New()
	testCases := []struct {
		name        string
		expectedID  uuid.UUID
		expectedErr error
		mockRepoFn  func(*mock.MockRepository)
	}{
		{
			name:        "Successful TestUsecase_CreateAccount",
			expectedID:  testUUID,
			expectedErr: nil,
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CreateAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(testUUID, nil)
			},
		},
		{
			name:        "Error in TestUsecase_CreateAccount",
			expectedID:  uuid.UUID{},
			expectedErr: fmt.Errorf("[usecase] can't create account into repository: some error"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CreateAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.UUID{}, errors.New("some error"))
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
			account := &models.Accounts{} // You should create an instance of your account model here

			accountID, err := mockUsecase.CreateAccount(context.Background(), userID, account)

			assert.Equal(t, tc.expectedID, accountID)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_UpdateAccount(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr error
		mockRepoFn  func(*mock.MockRepository)
	}{
		{
			name:        "Successful TestUsecase_UpdateAccount",
			expectedErr: nil,
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckForbidden(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepository.EXPECT().UpdateAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:        "Forbidden Error in TestUsecase_UpdateAccount",
			expectedErr: fmt.Errorf("[usecase] can't be update by user: some forbidden error"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckForbidden(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("some forbidden error"))
			},
		},
		{
			name:        "Update Error in TestUsecase_UpdateAccount",
			expectedErr: fmt.Errorf("[usecase] can't update account into repository: some update error"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckForbidden(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepository.EXPECT().UpdateAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("some update error"))
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
			account := &models.Accounts{ID: uuid.New()} // Assuming ID is a required field for account

			err := mockUsecase.UpdateAccount(context.Background(), userID, account)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_DeleteAccount(t *testing.T) {
	userIDTest := uuid.New()
	accountIDTest := uuid.New()

	testCases := []struct {
		name        string
		expectedErr error
		mockRepoFn  func(*mock.MockRepository)
	}{
		{
			name:        "Successful deletion",
			expectedErr: nil,
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckForbidden(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepository.EXPECT().DeleteAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:        "Error in CheckForbidden",
			expectedErr: fmt.Errorf("[usecase] can't be delete by user: %w", errors.New("forbidden")),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckForbidden(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("forbidden"))
			},
		},
		{
			name:        "Error in DeleteAccount",
			expectedErr: fmt.Errorf("[usecase] can't delete account into repository: %w", errors.New("repository error")),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckForbidden(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepository.EXPECT().DeleteAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("repository error"))
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

			err := mockUsecase.DeleteAccount(context.Background(), userIDTest, accountIDTest)

			if (tc.expectedErr == nil && err != nil) ||
				(tc.expectedErr != nil && err == nil) ||
				(tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}
