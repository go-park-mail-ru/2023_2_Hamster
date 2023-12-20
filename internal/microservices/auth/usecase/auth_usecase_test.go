package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth"
	mock "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/mocks"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_SignUp(t *testing.T) {
	userIdTest := uuid.New()
	testCases := []struct {
		name             string
		expectedUserID   uuid.UUID
		expectedLogin    string
		expectedUsername string
		expectedErr      error
		mockRepoFn       func(*mock.MockRepository)
	}{
		{
			name:             "Successful SignUp",
			expectedUserID:   userIdTest,
			expectedLogin:    "testLogin",
			expectedUsername: "testUser",
			expectedErr:      nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckLoginUnique(gomock.Any(), gomock.Any()).Return(true, nil)
				mockRepositry.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(userIdTest, nil)
			},
		},
		{
			name:        "Error in CheckLoginUnique",
			expectedErr: fmt.Errorf("[usecase] error checking login uniqueness: some error"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckLoginUnique(gomock.Any(), gomock.Any()).Return(false, errors.New("some error"))
			},
		},
		{
			name:        "Username Already Exists",
			expectedErr: fmt.Errorf("[usecase] user already exists"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckLoginUnique(gomock.Any(), gomock.Any()).Return(false, nil)
			},
		},
		{
			name:        "Error in CreateUser",
			expectedErr: fmt.Errorf("[usecase] cannot create user: some error"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckLoginUnique(gomock.Any(), gomock.Any()).Return(true, nil)
				mockRepositry.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(uuid.Nil, errors.New("some error"))
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

			input := auth.SignUpInput{
				Login:          "testLogin",
				Username:       "testUser",
				PlaintPassword: "testPassword",
			}

			userID, login, username, err := mockUsecase.SignUp(context.Background(), input)
			assert.Equal(t, tc.expectedUserID, userID)
			assert.Equal(t, tc.expectedLogin, login)
			assert.Equal(t, tc.expectedUsername, username)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_Login(t *testing.T) {
	userIdTest := uuid.New()
	testCases := []struct {
		name             string
		expectedUserID   uuid.UUID
		expectedLogin    string
		expectedUsername string
		expectedErr      error
		mockRepoFn       func(*mock.MockRepository)
	}{
		{
			name:             "Successful Login",
			expectedUserID:   userIdTest,
			expectedLogin:    "testLogin",
			expectedUsername: "testUser",
			expectedErr:      nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				user := &models.User{
					ID:       userIdTest,
					Login:    "testLogin",
					Password: "$argon2id$v=19$m=65536,t=1,p=4$YcsUni+F/VK3Vsjuw7Hb/Q$eZwr1GdO2/bkDRa2PGfTw5l6LPymywRdE9St5ot2Gv8",
					Username: "testUser",
				}
				mockRepositry.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any()).Return(user, nil)
			},
		},
		{
			name:           "Error in GetUserByLogin",
			expectedErr:    fmt.Errorf("[usecase] can't find user: some error"),
			expectedUserID: uuid.Nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any()).Return(&models.User{}, errors.New("some error"))
			},
		},
		{
			name:           "Incorrect Password",
			expectedErr:    fmt.Errorf("[usecase] password hash doesn't match the real one: %w", &models.IncorrectPasswordError{UserID: userIdTest}),
			expectedUserID: uuid.Nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				user := &models.User{
					ID:       userIdTest,
					Login:    "testLogin",
					Password: "$argon2id$v=19$m=65536,t=1,p=4$YcsUni+F/VK3Vsjuw7Hb/Q$eZwr1GdO2/bkDRa2PGfTw5l6LPymywRdE9St5ot2",
					Username: "testUser",
				}
				mockRepositry.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any()).Return(user, nil)
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

			userID, login, username, err := mockUsecase.Login(context.Background(), "testLogin", "testPassword")
			assert.Equal(t, tc.expectedUserID, userID)
			assert.Equal(t, tc.expectedLogin, login)
			assert.Equal(t, tc.expectedUsername, username)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_CheckLoginUnique(t *testing.T) {
	testCases := []struct {
		name             string
		expectedIsUnique bool
		expectedErr      error
		mockRepoFn       func(*mock.MockRepository)
	}{
		{
			name:             "Login is Unique",
			expectedIsUnique: true,
			expectedErr:      nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckLoginUnique(gomock.Any(), gomock.Any()).Return(true, nil)
			},
		},
		{
			name:             "Login is Not Unique",
			expectedIsUnique: false,
			expectedErr:      nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckLoginUnique(gomock.Any(), gomock.Any()).Return(false, nil)
			},
		},
		{
			name:             "Error in CheckLoginUnique",
			expectedIsUnique: false,
			expectedErr:      fmt.Errorf("[usecase] can't login unique check"),
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().CheckLoginUnique(gomock.Any(), gomock.Any()).Return(false, errors.New("some error"))
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

			isUnique, err := mockUsecase.CheckLoginUnique(context.Background(), "testLogin")
			assert.Equal(t, tc.expectedIsUnique, isUnique)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_GetByID(t *testing.T) {
	userIdTest := uuid.New()
	testCases := []struct {
		name         string
		inputUserID  uuid.UUID
		expectedUser *models.User
		expectedErr  error
		mockRepoFn   func(*mock.MockRepository)
	}{
		{
			name:        "Successful GetByID",
			inputUserID: userIdTest,
			expectedUser: &models.User{
				ID: userIdTest,
				// Add other expected user fields here
			},
			expectedErr: nil,
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetByID(gomock.Any(), userIdTest).Return(&models.User{
					ID: userIdTest,
					// Add other mocked user fields here
				}, nil)
			},
		},
		{
			name:         "Error in GetByID",
			inputUserID:  userIdTest,
			expectedUser: nil,
			expectedErr:  fmt.Errorf("[usecase] can't get user from repository some error"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetByID(gomock.Any(), userIdTest).Return(nil, errors.New("some error"))
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

			user, err := mockUsecase.GetByID(context.Background(), tc.inputUserID)

			assert.Equal(t, tc.expectedUser, user)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_ChangePassword(t *testing.T) {
	userIdTest := uuid.New()
	testCases := []struct {
		name        string
		input       auth.ChangePasswordInput
		expectedErr error
		mockRepoFn  func(*mock.MockRepository)
	}{
		{
			name: "Successful ChangePassword",
			input: auth.ChangePasswordInput{
				Login:       "testLogin",
				OldPassword: "testPassword",
				NewPassword: "newPassword",
			},
			expectedErr: nil,
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetUserByLogin(gomock.Any(), "testLogin").Return(&models.User{
					ID:       userIdTest,
					Password: "$argon2id$v=19$m=65536,t=1,p=4$YcsUni+F/VK3Vsjuw7Hb/Q$eZwr1GdO2/bkDRa2PGfTw5l6LPymywRdE9St5ot2Gv8", // Replace with actual hashed password
				}, nil)
				mockRepository.EXPECT().ChangePassword(gomock.Any(), userIdTest, gomock.Any()).Return(nil)
			},
		},
		{
			name: "Error in GetUserByLogin",
			input: auth.ChangePasswordInput{
				Login:       "testLogin",
				OldPassword: "oldPassword",
				NewPassword: "newPassword",
			},
			expectedErr: fmt.Errorf("[usecase] can't find user: some error"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetUserByLogin(gomock.Any(), "testLogin").Return(nil, errors.New("some error"))
			},
		},
		// Add more test cases for different scenarios (e.g., incorrect old password, hasher errors, repo errors)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			tc.mockRepoFn(mockRepo)

			mockUsecase := NewUsecase(mockRepo, *logger.NewLogger(context.TODO()))

			err := mockUsecase.ChangePassword(context.Background(), tc.input)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}
