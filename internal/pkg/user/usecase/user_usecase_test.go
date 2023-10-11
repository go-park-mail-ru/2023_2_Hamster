package usecase

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/mocks"
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
			mockRepoFn: func(mockRepositry *mocks.MockRepository) {
				mockRepositry.EXPECT().GetUserBalance(gomock.Any()).Return(100.0, nil)
			},
		},
		{
			name:            "Error in balance retrieval",
			expectedBalance: 0,
			expectedErr:     fmt.Errorf("[usecase] cant't get balance from repository fsdaffd"),
			mockRepoFn: func(mockRepositry *mocks.MockRepository) {
				mockRepositry.EXPECT().GetUserBalance(gomock.Any()).Return(0.0, errors.New("fsdaffd"))
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
			if !reflect.DeepEqual(tc.expectedErr, err) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}
