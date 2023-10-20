package usecase

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/delivery/http/transfer_models"
	mock "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_GetUserBalance(t *testing.T) {
	testCases := []struct {
		name            string
		expectedBalance float64
		expectedErr     error
		mockRepoFn      func(*mock.MockRepository)
	}{
		{
			name:            "Successful balance retrieval",
			expectedBalance: 100.0,
			expectedErr:     nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().GetUserBalance(gomock.Any()).Return(100.0, nil)
			},
		},
		{
			name:            "Error in balance retrieval",
			expectedBalance: 0,
			expectedErr:     fmt.Errorf("[usecase] can't get balance from repository some error"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().GetUserBalance(gomock.Any()).Return(0.0, errors.New("some error"))
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

			balance, err := mockUsecase.GetUserBalance(userID)

			assert.Equal(t, tc.expectedBalance, balance)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_GetPlannedBudget(t *testing.T) {
	testCases := []struct {
		name           string
		expectedBudget float64
		expectedErr    error
		mockRepoFn     func(*mock.MockRepository)
	}{
		{
			name:           "Successful budget retrieval",
			expectedBudget: 200.0,
			expectedErr:    nil,
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetPlannedBudget(gomock.Any()).Return(200.0, nil)
			},
		},
		{
			name:           "Error in budget retrieval",
			expectedBudget: 0,
			expectedErr:    fmt.Errorf("[usecase] can't get planned budget from repository some error"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetPlannedBudget(gomock.Any()).Return(0.0, errors.New("some error"))
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

			budget, err := mockUsecase.GetPlannedBudget(userID)

			assert.Equal(t, tc.expectedBudget, budget)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_GetAccounts(t *testing.T) {
	uuidTest := uuid.New()
	testCases := []struct {
		name             string
		expectedAccounts []models.Accounts
		expectedErr      error
		mockRepoFn       func(*mock.MockRepository)
	}{
		{
			name: "Successful accounts retrieval",
			expectedAccounts: []models.Accounts{
				{ID: uuidTest, UserID: uuidTest, Balance: 100.0, MeanPayment: "Account1"},
				{ID: uuidTest, UserID: uuidTest, Balance: 200.0, MeanPayment: "Account2"},
			},
			expectedErr: nil,
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetAccounts(gomock.Any()).Return([]models.Accounts{
					{ID: uuidTest, UserID: uuidTest, Balance: 100.0, MeanPayment: "Account1"},
					{ID: uuidTest, UserID: uuidTest, Balance: 200.0, MeanPayment: "Account2"}}, nil)
			},
		},
		{
			name:             "Error in accounts retrieval",
			expectedAccounts: nil,
			expectedErr:      fmt.Errorf("[usecase] can't get accounts from repository some error"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetAccounts(gomock.Any()).Return(nil, errors.New("some error"))
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

			accounts, err := mockUsecase.GetAccounts(userID)

			assert.Equal(t, tc.expectedAccounts, accounts)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_GetCurrentBudget(t *testing.T) {
	testCases := []struct {
		name                  string
		expectedCurrentBudget float64
		expectedErr           error
		mockRepoFn            func(*mock.MockRepository)
	}{
		{
			name:                  "Successful current budget",
			expectedCurrentBudget: 0.0,
			expectedErr:           nil,
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetPlannedBudget(gomock.Any()).Return(1700.0, nil)
				mockRepository.EXPECT().GetCurrentBudget(gomock.Any()).Return(1700.0, nil)
			},
		},
		{
			name:                  "Error in planned issue",
			expectedCurrentBudget: 0.0,
			expectedErr:           fmt.Errorf("[usecase] can't get planned budget from repository some error"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetPlannedBudget(gomock.Any()).Return(0.0, errors.New("some error"))
				mockRepository.EXPECT().GetCurrentBudget(gomock.Any()).Return(0.0, nil)
			},
		},
		{
			name:                  "Error in planned issue",
			expectedCurrentBudget: 0.0,
			expectedErr:           fmt.Errorf("[usecase] can't get current budget from repository some error"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetCurrentBudget(gomock.Any()).Return(0.0, errors.New("some error"))
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

			currentBudget, err := mockUsecase.GetCurrentBudget(userID)

			assert.Equal(t, tc.expectedCurrentBudget, currentBudget)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_GetUser(t *testing.T) {
	testUserID := uuid.New()
	testCases := []struct {
		name         string
		expectedUser *models.User
		expectedErr  error
		mockRepoFn   func(*mock.MockRepository)
	}{
		{
			name: "Success GetUser",
			expectedUser: &models.User{ID: testUserID,
				Username:      "kossmatoff",
				PlannedBudget: 100.0,
				Password:      "hash",
				AvatarURL:     "img/image.png",
				Salt:          "a"},

			expectedErr: fmt.Errorf("[usecase] can't get user from repository some error"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetByID(gomock.Any()).Return(&models.User{ID: testUserID,
					Username:      "kossmatoff",
					PlannedBudget: 100.0,
					Password:      "hash",
					AvatarURL:     "img/image.png",
					Salt:          "a"}, errors.New("some error"))
			},
		},
		{
			name:         "Error in UserGet issue",
			expectedUser: &models.User{},
			expectedErr:  nil,
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				user := &models.User{}
				mockRepository.EXPECT().GetByID(gomock.Any()).Return(user, nil)
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

			userActual, err := mockUsecase.GetUser(userID)

			assert.Equal(t, tc.expectedUser, userActual)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_GetFeed(t *testing.T) {
	testUserID := uuid.New()
	testCases := []struct {
		name         string
		expectedFeed *transfer_models.UserFeed
		expectedErr  error
		mockRepoFn   func(*mock.MockRepository)
	}{
		{
			name:         "Error getUser balance retrieval",
			expectedFeed: &transfer_models.UserFeed{},
			expectedErr:  fmt.Errorf("[usecase] can't get balance from repository some erros"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetUserBalance(gomock.Any()).Return(0.0, errors.New("some erros"))
			},
		},
		{
			name:         "Error getUser current budget retrieval",
			expectedFeed: &transfer_models.UserFeed{},
			expectedErr:  fmt.Errorf("[usecase] can't get current budget from repository some error"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetUserBalance(gomock.Any()).Return(0.0, nil)
				mockRepository.EXPECT().GetCurrentBudget(gomock.Any()).Return(0.0, errors.New("some error"))
			},
		},
		{
			name:         "Error in getUser planned budget retrieval",
			expectedFeed: &transfer_models.UserFeed{},
			expectedErr:  fmt.Errorf("[usecase] can't get planned budget from repository some error"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetUserBalance(gomock.Any()).Return(0.0, nil)
				mockRepository.EXPECT().GetCurrentBudget(gomock.Any()).Return(0.0, nil)
				mockRepository.EXPECT().GetPlannedBudget(gomock.Any()).Return(0.0, errors.New("some error"))
			},
		},
		// {
		// 	name:         "Error in getUser accounts retrieval",
		// 	expectedFeed: &transfer_models.UserFeed{},
		// 	expectedErr:  fmt.Errorf("[usecase] can't get accounts from repository some error"),
		// 	mockRepoFn: func(mockRepository *mock.MockRepository) {
		// 		mockRepository.EXPECT().GetUserBalance(testUserID).Return(0.0, nil)
		// 		mockRepository.EXPECT().GetCurrentBudget(testUserID).Return(0.0, nil)
		// 		mockRepository.EXPECT().GetPlannedBudget(testUserID).Return(0.0, nil)
		// 		mockRepository.EXPECT().GetAccounts(testUserID).Return([]models.Accounts{}, errors.New("some error"))
		// 	},
		// },
		// {
		// 	name:         "Success in getUser",
		// 	expectedFeed: &transfer_models.UserFeed{},

		// 	expectedErr: nil,
		// 	mockRepoFn: func(mockRepository *mock.MockRepository) {
		// 		mockRepository.EXPECT().GetUserBalance(gomock.Any()).Return(100.0, nil)
		// 		mockRepository.EXPECT().GetCurrentBudget(gomock.Any()).Return(100.0, nil)
		// 		mockRepository.EXPECT().GetPlannedBudget(gomock.Any()).Return(0.0, nil)
		// 		mockRepository.EXPECT().GetAccounts(gomock.Any()).Return([]models.Accounts{}, nil)
		// 	},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			tc.mockRepoFn(mockRepo)

			mockUsecase := NewUsecase(mockRepo, *logger.CreateCustomLogger())

			feedActual, err := mockUsecase.GetFeed(testUserID)

			assert.Equal(t, tc.expectedFeed, feedActual)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}