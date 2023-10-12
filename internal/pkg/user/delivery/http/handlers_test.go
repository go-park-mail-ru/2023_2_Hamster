package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	mocks "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/mocks"
	"github.com/gorilla/mux"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetUserBalance(t *testing.T) {
	uuidTest := uuid.New()
	tests := []struct {
		name          string
		userID        string
		expectedCode  int
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
	}{
		{
			name:         "Successful call to GetUserBalance",
			userID:       uuid.New().String(),
			expectedCode: http.StatusOK,
			expectedBody: `{"balance":100}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				expectedBalance := 100.0
				mockUsecase.EXPECT().GetUserBalance(gomock.Any()).Return(expectedBalance, nil)
			},
		},
		{
			name:         "Invalid userID",
			userID:       "invalidUserID",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"invalid uuid parameter"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				// No expectations for mockUsecase.
			},
		},
		{
			name:         "Error from GetUserBalance",
			userID:       uuidTest.String(),
			expectedCode: http.StatusBadRequest,
			expectedBody: fmt.Sprintf(`{"message":"balance from user: %s doesn't exist"}`, uuidTest),
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				errorUserID := uuidTest
				expectedError := models.NoSuchUserIdBalanceError{UserID: errorUserID}
				mockUsecase.EXPECT().GetUserBalance(gomock.Any()).Return(0.0, &expectedError)
			},
		},
		{
			name:         "Internal server error",
			userID:       uuid.New().String(),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"message":"internal server error"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				//internalErrorUserID := uuid.New()
				internalError := errors.New("internal server error")
				mockUsecase.EXPECT().GetUserBalance(gomock.Any()).Return(0.0, internalError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			mockHandler := NewHandler(mockUsecase, *logger.CreateCustomLogger())

			url := "/api/user/" + tt.userID + "/balance"
			req := httptest.NewRequest("GET", url, nil)
			req = mux.SetURLVars(req, map[string]string{"userID": tt.userID})

			recorder := httptest.NewRecorder()
			mockHandler.GetUserBalance(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_GetPlannedBudget(t *testing.T) {
	uuidTest := uuid.New()
	tests := []struct {
		name          string
		userID        string
		expectedCode  int
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
	}{
		{
			name:         "Successful call to GetPlannedBudget",
			userID:       uuid.New().String(),
			expectedCode: http.StatusOK,
			expectedBody: `{"planned_balance":100}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				expectedBudget := 100.0
				mockUsecase.EXPECT().GetPlannedBudget(gomock.Any()).Return(expectedBudget, nil)
			},
		},
		{
			name:         "Invalid userID",
			userID:       "invalidUserID",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"invalid uuid parameter"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				// No expectations for mockUsecase.
			},
		},
		{
			name:         "Error from GetPlannedBudget",
			userID:       uuidTest.String(),
			expectedCode: http.StatusBadRequest,
			expectedBody: fmt.Sprintf(`{"message":"planned budget from user: %s doesn't exist"}`, uuidTest),
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				errorUserID := uuidTest
				expectedError := models.NoSuchPlannedBudgetError{UserID: errorUserID}
				mockUsecase.EXPECT().GetPlannedBudget(gomock.Any()).Return(0.0, &expectedError)
			},
		},
		{
			name:         "Internal server error",
			userID:       uuid.New().String(),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"message":"internal server error"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				internalError := errors.New("internal server error")
				mockUsecase.EXPECT().GetPlannedBudget(gomock.Any()).Return(0.0, internalError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			mockHandler := NewHandler(mockUsecase, *logger.CreateCustomLogger())

			url := "/api/user/" + tt.userID + "/planned_budget"
			req := httptest.NewRequest("GET", url, nil)
			req = mux.SetURLVars(req, map[string]string{"userID": tt.userID})

			recorder := httptest.NewRecorder()
			mockHandler.GetPlannedBudget(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_GetCurrentBudget(t *testing.T) {
	uuidTest := uuid.New()
	tests := []struct {
		name          string
		userID        string
		expectedCode  int
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
	}{
		{
			name:         "Successful call to GetCurrentBudget",
			userID:       uuid.New().String(),
			expectedCode: http.StatusOK,
			expectedBody: `{"actual_balance":100}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				expectedBudget := 100.0
				mockUsecase.EXPECT().GetCurrentBudget(gomock.Any()).Return(expectedBudget, nil)
			},
		},
		{
			name:         "Invalid userID",
			userID:       "invalidUserID",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"invalid uuid parameter"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				// No expectations for mockUsecase.
			},
		},
		{
			name:         "Error from GetCurrentBudget",
			userID:       uuidTest.String(),
			expectedCode: http.StatusBadRequest,
			expectedBody: fmt.Sprintf(`{"message":"actual budget from user: %s doesn't exist"}`, uuidTest),
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				errorUserID := uuidTest
				expectedError := models.NoSuchCurrentBudget{UserID: errorUserID}
				mockUsecase.EXPECT().GetCurrentBudget(gomock.Any()).Return(0.0, &expectedError)
			},
		},
		{
			name:         "Internal server error",
			userID:       uuid.New().String(),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"message":"internal server error"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				internalError := errors.New("internal server error")
				mockUsecase.EXPECT().GetCurrentBudget(gomock.Any()).Return(0.0, internalError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			mockHandler := NewHandler(mockUsecase, *logger.CreateCustomLogger())

			url := "/api/user/" + tt.userID + "/current_budget"
			req := httptest.NewRequest("GET", url, nil)
			req = mux.SetURLVars(req, map[string]string{"userID": tt.userID})

			recorder := httptest.NewRecorder()
			mockHandler.GetCurrentBudget(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_GetAccounts(t *testing.T) {
	uuidTest := uuid.New()
	tests := []struct {
		name          string
		userID        string
		expectedCode  int
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
	}{
		{
			name:         "Successful call to GetAccounts",
			userID:       uuid.New().String(),
			expectedCode: http.StatusOK,
			expectedBody: `{"Account":[]}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().GetAccounts(gomock.Any()).Return([]models.Accounts{}, nil)
			},
		},
		{
			name:         "Invalid userID",
			userID:       "invalidUserID",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"invalid uuid parameter"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				// No expectations for mockUsecase.
			},
		},
		{
			name:         "No accounts found",
			userID:       uuidTest.String(),
			expectedCode: http.StatusOK,
			expectedBody: `""`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				expectedError := models.NoSuchAccounts{}
				mockUsecase.EXPECT().GetAccounts(gomock.Any()).Return([]models.Accounts{}, &expectedError)
			},
		},
		{
			name:         "Internal server error",
			userID:       uuid.New().String(),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"message":"internal server error"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				internalError := errors.New("internal server error")
				mockUsecase.EXPECT().GetAccounts(gomock.Any()).Return([]models.Accounts{}, internalError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			mockHandler := NewHandler(mockUsecase, *logger.CreateCustomLogger())

			url := "/api/user/" + tt.userID + "/accounts"
			req := httptest.NewRequest("GET", url, nil)
			req = mux.SetURLVars(req, map[string]string{"userID": tt.userID})

			recorder := httptest.NewRecorder()
			mockHandler.GetAccounts(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}
