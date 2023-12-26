package http

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	genCategory "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/delivery/grpc/generated"
	mocks "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/mocks"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateTag(t *testing.T) {
	uuidTest := uuid.New()
	user := &models.User{ID: uuidTest}
	tests := []struct {
		name           string
		user           *models.User
		expectedCode   int
		expectedBody   string
		mockUsecaseFn  func(usecase *mocks.MockCategoryServiceClient)
		requestPayload string
	}{
		{
			name:         "Successful Tag Creation",
			user:         user,
			expectedCode: http.StatusOK,
			expectedBody: `{"status":200,"body":{"category_id":"` + uuidTest.String() + `"}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockCategoryServiceClient) {
				mockUsecase.EXPECT().CreateTag(gomock.Any(), gomock.Any()).Return(&genCategory.CreateTagResponse{TagId: uuidTest.String()}, nil)
			},
			requestPayload: `{"name": "TestTag", "showIncome": true, "showOutcome": true, "regular": true}`,
		},
		{
			name:           "Unauthorized Request",
			user:           nil,
			expectedCode:   http.StatusUnauthorized,
			expectedBody:   `{"status":401,"message":"unauthorized"}`,
			mockUsecaseFn:  func(mockUsecase *mocks.MockCategoryServiceClient) {},
			requestPayload: `{"name": "TestTag", "showIncome": true, "showOutcome": true, "regular": true}`,
		},
		{
			name:           "Invalid Request Body",
			user:           user,
			expectedCode:   http.StatusBadRequest,
			expectedBody:   `{"status":400,"message":"Corrupted request body can't unmarshal"}`,
			mockUsecaseFn:  func(mockUsecase *mocks.MockCategoryServiceClient) {},
			requestPayload: `invalid_json`,
		},
		{
			name:         "Error in Tag Creation",
			user:         user,
			expectedCode: http.StatusTooManyRequests,
			expectedBody: `{"status":429,"message":"Can't create tag"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockCategoryServiceClient) {
				mockUsecase.EXPECT().CreateTag(gomock.Any(), gomock.Any()).Return(nil, errors.New("tag not created"))
			},
			requestPayload: `{"name": "TestTag", "showIncome": true, "showOutcome": true, "regular": true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockCategoryServiceClient(ctrl)
			tt.mockUsecaseFn(mockService)

			mockHandler := NewHandler(mockService, *logger.NewLogger(context.TODO()))

			req := httptest.NewRequest("POST", "/api/tag/create", strings.NewReader(tt.requestPayload))
			req.Header.Set("Content-Type", "application/json")

			if tt.user != nil {
				ctx := context.WithValue(req.Context(), models.ContextKeyUserType{}, tt.user)
				req = req.WithContext(ctx)
			}

			recorder := httptest.NewRecorder()

			mockHandler.CreateTag(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

// func TestHandler_GetTags(t *testing.T) {
// 	uuidTest := uuid.New()
// 	user := &models.User{ID: uuidTest}
// 	tests := []struct {
// 		name          string
// 		user          *models.User
// 		expectedCode  int
// 		expectedBody  string
// 		mockUsecaseFn func(mockUsecase *mocks.MockCategoryServiceClient)
// 	}{
// 		{
// 			name:         "Successful Get Tags",
// 			user:         user,
// 			expectedCode: http.StatusOK,
// 			expectedBody: `{"status":200,"body":[{"id":"` + uuidTest.String() + `","user_id":"` + uuidTest.String() + `","parent_id":"` + uuid.Nil.String() + `","name":"TestTag","show_income":true,"show_outcome":true,"regular":true}]}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockCategoryServiceClient) {
// 				mockUsecase.EXPECT().GetTags(gomock.Any(), gomock.Any()).Return(&genCategory.GetTagsResponse{
// 					Categories: []*genCategory.Category{
// 						{
// 							Id:          uuidTest.String(),
// 							UserId:      uuidTest.String(),
// 							ParentId:    "",
// 							Name:        "TestTag",
// 							ShowIncome:  true,
// 							ShowOutcome: true,
// 							Regular:     true,
// 						},
// 					},
// 				}, nil)
// 			},
// 		},
// 		{
// 			name:          "Unauthorized Request",
// 			user:          nil,
// 			expectedCode:  http.StatusUnauthorized,
// 			expectedBody:  `{"status":401,"message":"unauthorized"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockCategoryServiceClient) {},
// 		},
// 		{
// 			name:         "Error in Get Tags",
// 			user:         user,
// 			expectedCode: http.StatusTooManyRequests,
// 			expectedBody: `{"status":429,"message":"Can't get tags"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockCategoryServiceClient) {
// 				mockUsecase.EXPECT().GetTags(gomock.Any(), gomock.Any()).Return(nil, errors.New("error getting tags"))
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockService := mocks.NewMockCategoryServiceClient(ctrl)
// 			tt.mockUsecaseFn(mockService)

// 			mockHandler := NewHandler(mockService, *logger.NewLogger(context.TODO()))

// 			req := httptest.NewRequest("GET", "/api/tags", nil)

// 			if tt.user != nil {
// 				ctx := context.WithValue(req.Context(), models.ContextKeyUserType{}, tt.user)
// 				req = req.WithContext(ctx)
// 			}

// 			recorder := httptest.NewRecorder()

// 			mockHandler.GetTags(recorder, req)

// 			actual := strings.TrimSpace(recorder.Body.String())

// 			assert.Equal(t, tt.expectedCode, recorder.Code)
// 			assert.Equal(t, tt.expectedBody, actual)
// 		})
// 	}
// }

func TestHandler_UpdateTag(t *testing.T) {
	uuidTest := uuid.New()
	user := &models.User{ID: uuidTest}
	tagID := uuid.New()
	tag := models.Category{
		ID:          tagID,
		UserID:      uuidTest,
		ParentID:    uuid.Nil,
		Name:        "TestTag",
		ShowIncome:  true,
		ShowOutcome: true,
		Regular:     true,
	}

	tests := []struct {
		name           string
		user           *models.User
		tag            models.Category
		expectedCode   int
		expectedBody   string
		mockUsecaseFn  func(mockUsecase *mocks.MockCategoryServiceClient)
		requestPayload string
	}{
		{
			name:         "Successful Tag Update",
			user:         user,
			tag:          tag,
			expectedCode: http.StatusOK,
			expectedBody: `{"status":200,"body":{"id":"` + tagID.String() + `","user_id":"` + uuidTest.String() + `","parent_id":"` + uuid.Nil.String() + `","name":"TestTag","show_income":true,"show_outcome":true,"regular":true}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockCategoryServiceClient) {
				mockUsecase.EXPECT().UpdateTag(gomock.Any(), gomock.Any()).Return(&genCategory.Category{
					Id:          tagID.String(),
					UserId:      uuidTest.String(),
					ParentId:    "",
					Name:        "TestTag",
					ShowIncome:  true,
					ShowOutcome: true,
					Regular:     true,
				}, nil)
			},
			requestPayload: `{"name": "TestTag", "showIncome": true, "showOutcome": true, "regular": true}`,
		},
		{
			name:           "Unauthorized Request",
			user:           nil,
			tag:            tag,
			expectedCode:   http.StatusUnauthorized,
			expectedBody:   `{"status":401,"message":"unauthorized"}`,
			mockUsecaseFn:  func(mockUsecase *mocks.MockCategoryServiceClient) {},
			requestPayload: `{"name": "TestTag", "showIncome": true, "showOutcome": true, "regular": true}`,
		},
		{
			name:           "Invalid Request Body",
			user:           user,
			tag:            models.Category{},
			expectedCode:   http.StatusBadRequest,
			expectedBody:   `{"status":400,"message":"Corrupted request body can't unmarshal"}`,
			mockUsecaseFn:  func(mockUsecase *mocks.MockCategoryServiceClient) {},
			requestPayload: `{"invalid_json": "missing_quote}`, // Invalid JSON string
		},
		{
			name:         "Error in Tag Update",
			user:         user,
			tag:          tag,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"Can't Update tag"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockCategoryServiceClient) {
				mockUsecase.EXPECT().UpdateTag(gomock.Any(), gomock.Any()).Return(nil, errors.New("tag not updated"))
			},
			requestPayload: `{"name": "TestTag", "showIncome": true, "showOutcome": true, "regular": true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockCategoryServiceClient(ctrl)
			tt.mockUsecaseFn(mockService)

			mockHandler := NewHandler(mockService, *logger.NewLogger(context.TODO()))
			reqBody := []byte(tt.requestPayload)

			req := httptest.NewRequest("PUT", "/api/tag/update", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.user != nil {
				ctx := context.WithValue(req.Context(), models.ContextKeyUserType{}, tt.user)
				req = req.WithContext(ctx)
			}

			recorder := httptest.NewRecorder()

			mockHandler.UpdateTag(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_DeleteTag(t *testing.T) {
	uuidTest := uuid.New()
	user := &models.User{ID: uuidTest}
	tagID := uuid.New()

	tests := []struct {
		name           string
		user           *models.User
		tagID          uuid.UUID
		expectedCode   int
		expectedBody   string
		mockUsecaseFn  func(mockUsecase *mocks.MockCategoryServiceClient)
		requestPayload string
	}{
		{
			name:         "Successful Tag Deletion",
			user:         user,
			tagID:        tagID,
			expectedCode: http.StatusOK,
			expectedBody: `{"status":200,"body":{"id":"` + tagID.String() + `"}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockCategoryServiceClient) {
				mockUsecase.EXPECT().DeleteTag(gomock.Any(), &genCategory.DeleteRequest{
					TagId:  tagID.String(),
					UserId: user.ID.String(),
				}).Return(nil, nil)
			},
			requestPayload: `{"id": "` + tagID.String() + `"}`,
		},
		{
			name:           "Unauthorized Request",
			user:           nil,
			tagID:          tagID,
			expectedCode:   http.StatusUnauthorized,
			expectedBody:   `{"status":401,"message":"unauthorized"}`,
			mockUsecaseFn:  func(mockUsecase *mocks.MockCategoryServiceClient) {},
			requestPayload: `{"id": "` + tagID.String() + `"}`,
		},
		{
			name:           "Invalid Request Body",
			user:           user,
			tagID:          tagID,
			expectedCode:   http.StatusBadRequest,
			expectedBody:   `{"status":400,"message":"Corupted request body can't unmarshal"}`,
			mockUsecaseFn:  func(mockUsecase *mocks.MockCategoryServiceClient) {},
			requestPayload: `{"invalid_json": "missing_quote}`, // Invalid JSON string
		},
		{
			name:         "Error in Tag Deletion",
			user:         user,
			tagID:        tagID,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"Can't Delete tag"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockCategoryServiceClient) {
				mockUsecase.EXPECT().DeleteTag(gomock.Any(), &genCategory.DeleteRequest{
					TagId:  tagID.String(),
					UserId: user.ID.String(),
				}).Return(nil, errors.New("tag not deleted"))
			},
			requestPayload: `{"id": "` + tagID.String() + `"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockCategoryServiceClient(ctrl)
			tt.mockUsecaseFn(mockService)

			mockHandler := NewHandler(mockService, *logger.NewLogger(context.TODO()))
			reqBody := []byte(tt.requestPayload)

			req := httptest.NewRequest("DELETE", "/api/tag/delete", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.user != nil {
				ctx := context.WithValue(req.Context(), models.ContextKeyUserType{}, tt.user)
				req = req.WithContext(ctx)
			}

			recorder := httptest.NewRecorder()

			mockHandler.DeleteTag(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}
