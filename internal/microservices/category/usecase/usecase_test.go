package usecase

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category"
	mock "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/mocks"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestUsecase_CreateTag(t *testing.T) {
	testID := uuid.New()
	testCases := []struct {
		name        string
		expectedErr error
		expectedID  uuid.UUID
		mockRepoFn  func(*mock.MockRepository)
	}{
		{
			name:        "Successful create tag",
			expectedErr: nil,
			expectedID:  testID,
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckNameUniq(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				mockRepository.EXPECT().CreateTag(gomock.Any(), gomock.Any()).Return(testID, nil)
			},
		},
		{
			name:        "Error CheckNameUniq",
			expectedErr: fmt.Errorf("[usecase] Error check uniq uniqueness: %v", "some error"),
			expectedID:  uuid.Nil,
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckNameUniq(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, errors.New("some error"))
			},
		},
		{
			name:        "Error CreateTag",
			expectedErr: fmt.Errorf("[usecase] Error tag creation: %v", "some error"),
			expectedID:  uuid.Nil,
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckNameUniq(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				mockRepository.EXPECT().CreateTag(gomock.Any(), gomock.Any()).Return(uuid.Nil, errors.New("some error"))
			},
		},
		{
			name:        "Error category already exist",
			expectedErr: fmt.Errorf("[usecase] Error category alredy exist"),
			expectedID:  uuid.Nil,
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckNameUniq(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
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

			tagInput := category.TagInput{}

			id, err := mockUsecase.CreateTag(context.Background(), tagInput)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}

			if id != tc.expectedID {
				t.Errorf("Expected ID: %v, but got: %v", tc.expectedID, id)
			}
		})
	}
}

func TestUsecase_UpdateTag(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr error
		mockRepoFn  func(*mock.MockRepository)
	}{
		{
			name:        "Successful update tag",
			expectedErr: nil,
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				mockRepository.EXPECT().UpdateTag(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:        "Error in check tag existence",
			expectedErr: fmt.Errorf("[usecase] Error in check tag existens: %v", "some error"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, errors.New("some error"))
			},
		},
		{
			name:        "Error tag doesn't exist can't update",
			expectedErr: fmt.Errorf("[usecase] Error tag doesn't exist can't update"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
			},
		},
		{
			name:        "Update tag Error",
			expectedErr: fmt.Errorf("[usecase] update tag Error: %v", "some error"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				mockRepository.EXPECT().UpdateTag(gomock.Any(), gomock.Any()).Return(errors.New("some error"))
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

			tag := &models.Category{}

			err := mockUsecase.UpdateTag(context.Background(), tag)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_DeleteTag(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr error
		mockRepoFn  func(*mock.MockRepository)
	}{
		{
			name:        "Successful delete tag",
			expectedErr: nil,
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				mockRepository.EXPECT().DeleteTag(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:        "Error in check tag existence",
			expectedErr: fmt.Errorf("[usecase] Error in check tag existens: %v", "some error"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, errors.New("some error"))
			},
		},
		{
			name:        "Error tag doesn't exist can't delete",
			expectedErr: fmt.Errorf("[usecase] Error tag doesn't exist can't delete"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
			},
		},
		{
			name:        "Error in tag deletion",
			expectedErr: fmt.Errorf("[usecase] Error in tag deletion: %v", "some error"),
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().CheckExist(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				mockRepository.EXPECT().DeleteTag(gomock.Any(), gomock.Any()).Return(errors.New("some error"))
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

			tagId := uuid.New()
			userId := uuid.New()

			err := mockUsecase.DeleteTag(context.Background(), tagId, userId)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUsecase_GetTags(t *testing.T) {
	testCases := []struct {
		name         string
		expectedErr  error
		expectedTags []models.Category
		mockRepoFn   func(*mock.MockRepository)
	}{
		{
			name:         "Successful get tags",
			expectedErr:  nil,
			expectedTags: []models.Category{{}, {}},
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetTags(gomock.Any(), gomock.Any()).Return([]models.Category{{}, {}}, nil)
			},
		},
		{
			name:         "Error getting user tags",
			expectedErr:  fmt.Errorf("[usecase] Error getting user tags: %v", "some error"),
			expectedTags: nil,
			mockRepoFn: func(mockRepository *mock.MockRepository) {
				mockRepository.EXPECT().GetTags(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
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

			userId := uuid.New()

			tags, err := mockUsecase.GetTags(context.Background(), userId)

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}

			if !reflect.DeepEqual(tags, tc.expectedTags) {
				t.Errorf("Expected tags: %v, got: %v", tc.expectedTags, tags)
			}
		})
	}
}
