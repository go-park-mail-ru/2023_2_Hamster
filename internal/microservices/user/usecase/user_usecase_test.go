package usecase

//import (
//	"context"
//	"errors"
//	"fmt"
//	"testing"
//
//	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
//	mock_account "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account/mocks"
//	mock "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/mocks"
//	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
//	"github.com/golang/mock/gomock"
//	"github.com/google/uuid"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestUsecase_GetUserBalance(t *testing.T) {
//	testCases := []struct {
//		name            string
//		expectedBalance float64
//		expectedErr     error
//		mockRepoFn      func(*mock.MockRepository)
//	}{
//		{
//			name:            "Successful balance retrieval",
//			expectedBalance: 100.0,
//			expectedErr:     nil,
//			mockRepoFn: func(mockRepositry *mock.MockRepository) {
//				mockRepositry.EXPECT().GetUserBalance(gomock.Any(), gomock.Any()).Return(100.0, nil)
//			},
//		},
//		{
//			name:            "Error in balance retrieval",
//			expectedBalance: 0,
//			expectedErr:     fmt.Errorf("[usecase] can't get balance from repository some error"),
//			mockRepoFn: func(mockRepositry *mock.MockRepository) {
//				mockRepositry.EXPECT().GetUserBalance(gomock.Any(), gomock.Any()).Return(0.0, errors.New("some error"))
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			mockRepo := mock.NewMockRepository(ctrl)
//			tc.mockRepoFn(mockRepo)
//			mockRepoa := mock_account.NewMockRepository(ctrl)
//
//			mockUsecase := NewUsecase(mockRepo, *logger.NewLogger(context.TODO()), mockRepoa)
//
//			userID := uuid.New()
//
//			balance, err := mockUsecase.GetUserBalance(context.Background(), userID)
//
//			assert.Equal(t, tc.expectedBalance, balance)
//			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
//				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
//			}
//		})
//	}
//}
//
//func TestUsecase_GetPlannedBudget(t *testing.T) {
//	testCases := []struct {
//		name           string
//		expectedBudget float64
//		expectedErr    error
//		mockRepoFn     func(*mock.MockRepository)
//	}{
//		{
//			name:           "Successful budget retrieval",
//			expectedBudget: 200.0,
//			expectedErr:    nil,
//			mockRepoFn: func(mockRepository *mock.MockRepository) {
//				mockRepository.EXPECT().GetPlannedBudget(gomock.Any(), gomock.Any()).Return(200.0, nil)
//			},
//		},
//		{
//			name:           "Error in budget retrieval",
//			expectedBudget: 0,
//			expectedErr:    fmt.Errorf("[usecase] can't get planned budget from repository some error"),
//			mockRepoFn: func(mockRepository *mock.MockRepository) {
//				mockRepository.EXPECT().GetPlannedBudget(gomock.Any(), gomock.Any()).Return(0.0, errors.New("some error"))
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			mockRepo := mock.NewMockRepository(ctrl)
//			tc.mockRepoFn(mockRepo)
//			mockRepoa := mock_account.NewMockRepository(ctrl)
//
//			mockUsecase := NewUsecase(mockRepo, *logger.NewLogger(context.TODO()), mockRepoa)
//
//			userID := uuid.New()
//
//			budget, err := mockUsecase.GetPlannedBudget(context.Background(), userID)
//
//			assert.Equal(t, tc.expectedBudget, budget)
//			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
//				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
//			}
//		})
//	}
//}
//
//func TestUsecase_GetAccounts(t *testing.T) {
//	uuidTest := uuid.New()
//	testCases := []struct {
//		name             string
//		expectedAccounts []models.Accounts
//		expectedErr      error
//		mockRepoFn       func(*mock.MockRepository)
//	}{
//		{
//			name: "Successful accounts retrieval",
//			expectedAccounts: []models.Accounts{
//				{ID: uuidTest, Balance: 100.0, MeanPayment: "Account1"},
//				{ID: uuidTest, Balance: 200.0, MeanPayment: "Account2"},
//			},
//			expectedErr: nil,
//			mockRepoFn: func(mockRepository *mock.MockRepository) {
//				mockRepository.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).Return([]models.Accounts{
//					{ID: uuidTest, Balance: 100.0, MeanPayment: "Account1"},
//					{ID: uuidTest, Balance: 200.0, MeanPayment: "Account2"}}, nil)
//			},
//		},
//		{
//			name:             "Error in accounts retrieval",
//			expectedAccounts: nil,
//			expectedErr:      fmt.Errorf("[usecase] can't get accounts from repository some error"),
//			mockRepoFn: func(mockRepository *mock.MockRepository) {
//				mockRepository.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			mockRepo := mock.NewMockRepository(ctrl)
//			tc.mockRepoFn(mockRepo)
//			mockRepoa := mock_account.NewMockRepository(ctrl)
//
//			mockUsecase := NewUsecase(mockRepo, *logger.NewLogger(context.TODO()), mockRepoa)
//
//			userID := uuid.New()
//
//			accounts, err := mockUsecase.GetAccounts(context.Background(), userID)
//
//			assert.Equal(t, tc.expectedAccounts, accounts)
//			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
//				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
//			}
//		})
//	}
//}
//
//func TestUsecase_GetCurrentBudget(t *testing.T) {
//	testCases := []struct {
//		name                  string
//		expectedCurrentBudget float64
//		expectedErr           error
//		mockRepoFn            func(*mock.MockRepository)
//	}{
//		{
//			name:                  "Successful current budget",
//			expectedCurrentBudget: 0.0,
//			expectedErr:           nil,
//			mockRepoFn: func(mockRepository *mock.MockRepository) {
//				mockRepository.EXPECT().GetPlannedBudget(gomock.Any(), gomock.Any()).Return(1700.0, nil)
//				mockRepository.EXPECT().GetCurrentBudget(gomock.Any(), gomock.Any()).Return(1700.0, nil)
//			},
//		},
//		{
//			name:                  "Error in planned issue",
//			expectedCurrentBudget: 0.0,
//			expectedErr:           fmt.Errorf("[usecase] can't get planned budget from repository some error"),
//			mockRepoFn: func(mockRepository *mock.MockRepository) {
//				mockRepository.EXPECT().GetPlannedBudget(gomock.Any(), gomock.Any()).Return(0.0, errors.New("some error"))
//				mockRepository.EXPECT().GetCurrentBudget(gomock.Any(), gomock.Any()).Return(0.0, nil)
//			},
//		},
//		{
//			name:                  "Error in planned issue",
//			expectedCurrentBudget: 0.0,
//			expectedErr:           fmt.Errorf("[usecase] can't get current budget from repository some error"),
//			mockRepoFn: func(mockRepository *mock.MockRepository) {
//				mockRepository.EXPECT().GetCurrentBudget(gomock.Any(), gomock.Any()).Return(0.0, errors.New("some error"))
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			mockRepo := mock.NewMockRepository(ctrl)
//			tc.mockRepoFn(mockRepo)
//			mockRepoa := mock_account.NewMockRepository(ctrl)
//
//			mockUsecase := NewUsecase(mockRepo, *logger.NewLogger(context.TODO()), mockRepoa)
//
//			userID := uuid.New()
//
//			currentBudget, err := mockUsecase.GetCurrentBudget(context.Background(), userID)
//
//			assert.Equal(t, tc.expectedCurrentBudget, currentBudget)
//			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
//				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
//			}
//		})
//	}
//}
//
//// func TestUsecase_GetUser(t *testing.T) {
//// 	testUserID := uuid.New()
//// 	testCases := []struct {
//// 		name         string
//// 		expectedUser *models.User
//// 		expectedErr  error
//// 		mockRepoFn   func(*mock.MockRepository)
//// 	}{
//// 		{
//// 			name: "Success GetUser",
//// 			expectedUser: &models.User{ID: testUserID,
//// 				Username:      "kossmatoff",
//// 				PlannedBudget: 100.0,
//// 				Password:      "hash",
//// 				AvatarURL:     uuid.Nil,
//// 			},
//
//// 			expectedErr: fmt.Errorf("[usecase] can't get user from repository some error"),
//// 			mockRepoFn: func(mockRepository *mock.MockRepository) {
//// 				mockRepository.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&models.User{ID: testUserID,
//// 					Username:      "kossmatoff",
//// 					PlannedBudget: 100.0,
//// 					Password:      "hash",
//// 					AvatarURL:     uuid.Nil,
//// 				}, errors.New("some error"))
//// 			},
//// 		},
//// 		{
//// 			name:         "Error in UserGet issue",
//// 			expectedUser: &models.User{},
//// 			expectedErr:  nil,
//// 			mockRepoFn: func(mockRepository *mock.MockRepository) {
//// 				user := &models.User{}
//// 				mockRepository.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(user, nil)
//// 			},
//// 		},
//// 	}
//
//// 	for _, tc := range testCases {
//// 		t.Run(tc.name, func(t *testing.T) {
//// 			ctrl := gomock.NewController(t)
//// 			defer ctrl.Finish()
//
//// 			mockRepo := mock.NewMockRepository(ctrl)
//// 			tc.mockRepoFn(mockRepo)
//
////			mockUsecase := NewUsecase(mockRepo, *logger.NewLogger(context.TODO()))
//
//// 			userID := uuid.New()
//
//// 			userActual, err := mockUsecase.GetUser(context.Background(), userID)
//
//// 			assert.Equal(t, tc.expectedUser, userActual)
//// 			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
//// 				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
//// 			}
//// 		})
//// 	}
//// }
//
//func TestUsecase_GetFeed(t *testing.T) {
//	testUserID := uuid.New()
//	testCases := []struct {
//		name string
//		//expectedFeed *transfer_models.UserFeed
//		expectedErr error
//		mockRepoFn  func(*mock.MockRepository)
//	}{
//		{
//			name: "Error getUser balance retrieval",
//			//expectedFeed: &transfer_models.UserFeed{},
//			expectedErr: fmt.Errorf("[usecase] can't get balance from repository some erros"),
//			mockRepoFn: func(mockRepository *mock.MockRepository) {
//				mockRepository.EXPECT().GetUserBalance(gomock.Any(), gomock.Any()).Return(0.0, errors.New("some erros"))
//			},
//		},
//		{
//			name: "Error getUser account",
//			//expectedFeed: &transfer_models.UserFeed{Account: []models.Accounts{}},
//			expectedErr: fmt.Errorf("[usecase] can't get accounts from repository some error"),
//			mockRepoFn: func(mockRepository *mock.MockRepository) {
//				mockRepository.EXPECT().GetUserBalance(gomock.Any(), gomock.Any()).Return(0.0, nil)
//				mockRepository.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).Return([]models.Accounts{}, errors.New("some error"))
//
//			},
//		},
//		{
//			name: "Error in getUser planned budget",
//			//expectedFeed: &transfer_models.UserFeed{},
//			expectedErr: fmt.Errorf("[usecase] can't get planned budget from repository err"),
//			mockRepoFn: func(mockRepository *mock.MockRepository) {
//				mockRepository.EXPECT().GetUserBalance(gomock.Any(), gomock.Any()).Return(0.0, nil)
//				//mockRepository.EXPECT().GetPlannedBudget(gomock.Any(), gomock.Any()).Return(0.0, nil)
//				mockRepository.EXPECT().GetPlannedBudget(gomock.Any(), gomock.Any()).Return(0.0, errors.New("err"))
//				//mockRepository.EXPECT().GetCurrentBudget(gomock.Any(), gomock.Any()).Return(0.0, nil)
//				mockRepository.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).Return([]models.Accounts{}, nil)
//			},
//		},
//		{
//			name: "Error in getUser Planned budget retrieval",
//			//expectedFeed: &transfer_models.UserFeed{},
//			expectedErr: fmt.Errorf("[usecase] can't get current budget from repository err"),
//			mockRepoFn: func(mockRepository *mock.MockRepository) {
//				mockRepository.EXPECT().GetUserBalance(gomock.Any(), gomock.Any()).Return(0.0, nil)
//				mockRepository.EXPECT().GetPlannedBudget(gomock.Any(), gomock.Any()).Return(0.0, nil)
//				mockRepository.EXPECT().GetCurrentBudget(gomock.Any(), gomock.Any()).Return(0.0, errors.New("err"))
//				mockRepository.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).Return([]models.Accounts{}, nil)
//			},
//		},
//		{
//			name: "Successful",
//			//expectedFeed: &transfer_models.UserFeed{},
//			expectedErr: nil,
//			mockRepoFn: func(mockRepository *mock.MockRepository) {
//				mockRepository.EXPECT().GetUserBalance(gomock.Any(), gomock.Any()).Return(0.0, nil)
//				mockRepository.EXPECT().GetPlannedBudget(gomock.Any(), gomock.Any()).Return(0.0, nil)
//				mockRepository.EXPECT().GetPlannedBudget(gomock.Any(), gomock.Any()).Return(0.0, nil)
//				mockRepository.EXPECT().GetCurrentBudget(gomock.Any(), gomock.Any()).Return(0.0, nil)
//				mockRepository.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).Return([]models.Accounts{}, nil)
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			mockRepo := mock.NewMockRepository(ctrl)
//			tc.mockRepoFn(mockRepo)
//			mockRepoa := mock_account.NewMockRepository(ctrl)
//
//			mockUsecase := NewUsecase(mockRepo, *logger.NewLogger(context.TODO()), mockRepoa)
//
//			_, err := mockUsecase.GetFeed(context.Background(), testUserID)
//
//			//assert.Equal(t, tc.expectedFeed, feedActual)
//			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
//				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
//			}
//		})
//	}
//}
//
//func TestUsecase_UpdateUser(t *testing.T) {
//	testCases := []struct {
//		name        string
//		expectedErr error
//		mockRepoFn  func(*mock.MockRepository)
//	}{
//		{
//			name:        "Successful update",
//			expectedErr: nil,
//			mockRepoFn: func(mockRepository *mock.MockRepository) {
//				mockRepository.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(nil)
//			},
//		},
//		{
//			name:        "Error Update User",
//			expectedErr: fmt.Errorf("[usecase] can't update user some error"),
//			mockRepoFn: func(mockRepository *mock.MockRepository) {
//				mockRepository.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(errors.New("some error"))
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			mockRepo := mock.NewMockRepository(ctrl)
//			tc.mockRepoFn(mockRepo)
//			mockRepoa := mock_account.NewMockRepository(ctrl)
//
//			mockUsecase := NewUsecase(mockRepo, *logger.NewLogger(context.TODO()), mockRepoa)
//
//			user := &models.User{}
//
//			err := mockUsecase.UpdateUser(context.Background(), user)
//
//			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
//				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
//			}
//		})
//	}
//}
//
//func TestUsecase_UpdatePhoto(t *testing.T) {
//	testCases := []struct {
//		name        string
//		expectedErr error
//		mockRepoFn  func(*mock.MockRepository)
//	}{
//		{
//			name:        "Successful update",
//			expectedErr: nil,
//			mockRepoFn: func(mockRepository *mock.MockRepository) {
//				mockRepository.EXPECT().UpdatePhoto(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
//			},
//		},
//		{
//			name:        "Error",
//			expectedErr: fmt.Errorf("some err"),
//			mockRepoFn: func(mockRepository *mock.MockRepository) {
//				mockRepository.EXPECT().UpdatePhoto(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("some err"))
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			mockRepo := mock.NewMockRepository(ctrl)
//			tc.mockRepoFn(mockRepo)
//
//			mockRepoa := mock_account.NewMockRepository(ctrl)
//
//			mockUsecase := NewUsecase(mockRepo, *logger.NewLogger(context.TODO()), mockRepoa)
//
//			userID := uuid.New()
//
//			_, err := mockUsecase.UpdatePhoto(context.Background(), userID)
//
//			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
//				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
//			}
//		})
//	}
//}
