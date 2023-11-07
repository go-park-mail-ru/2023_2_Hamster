package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"

	mock "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/mocks"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_Login(t *testing.T) {
	userIdTest := uuid.New()
	testCases := []struct {
		name             string
		expectedUserID   uuid.UUID
		expectedUsername string
		expectedErr      error
		mockRepoFn       func(*mock.MockRepository)
	}{
		{
			name:             "Successful Login",
			expectedUserID:   userIdTest,
			expectedUsername: "testUser",
			expectedErr:      nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				user := models.User{
					ID:       userIdTest,
					Login:    "testLogin",
					Password: "hashedPassword",
					Username: "testUser",
				}
				mockRepositry.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any()).Return(user, nil)
				// mockRepositry.EXPECT().VerifyPassword(gomock.Any(), gomock.Any()).Return(true, nil)
			},
		},
		{
			name:           "Error in GetUserByLogin",
			expectedErr:    fmt.Errorf("[usecase] can't find user: some error"),
			expectedUserID: uuid.Nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				mockRepositry.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any()).Return(models.User{}, errors.New("some error"))
			},
		},
		{
			name:           "Incorrect Password",
			expectedErr:    fmt.Errorf("[usecase] incorrect password"),
			expectedUserID: uuid.Nil,
			mockRepoFn: func(mockRepositry *mock.MockRepository) {
				user := models.User{
					ID:       userIdTest,
					Login:    "testLogin",
					Password: "hashedPassword",
					Username: "testUser",
				}
				mockRepositry.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any()).Return(user, nil)
				// mockRepositry.EXPECT().VerifyPassword(gomock.Any(), gomock.Any()).Return(false, nil)
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

			userID, username, err := mockUsecase.Login(context.Background(), "testLogin", "testPassword")
			assert.Equal(t, tc.expectedUserID, userID)
			assert.Equal(t, tc.expectedUsername, username)
			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}
