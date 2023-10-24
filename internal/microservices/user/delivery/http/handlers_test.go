package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microsevices/user/delivery/http/transfer_models"
	mocks "github.com/go-park-mail-ru/2023_2_Hamster/internal/microsevices/user/mocks"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
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
			expectedBody: `{"status":200,"body":{"balance":100}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				expectedBalance := 100.0
				mockUsecase.EXPECT().GetUserBalance(gomock.Any()).Return(expectedBalance, nil)
			},
		},
		{
			name:         "Invalid userID",
			userID:       "invalidUserID",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"invalid url parameter"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				// No expectations for mockUsecase.
			},
		},
		{
			name:         "Error from GetUserBalance",
			userID:       uuidTest.String(),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"no such balance"}`,
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
			expectedBody: `{"status":500,"message":"can't get balance"}`,
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
			expectedBody: `{"status":200,"body":{"planned_balance":100}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				expectedBudget := 100.0
				mockUsecase.EXPECT().GetPlannedBudget(gomock.Any()).Return(expectedBudget, nil)
			},
		},
		{
			name:         "Invalid userID",
			userID:       "invalidUserID",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"invalid url parameter"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				// No expectations for mockUsecase.
			},
		},
		{
			name:         "Error from GetPlannedBudget",
			userID:       uuidTest.String(),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"no such planned budget"}`,
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
			expectedBody: `{"status":500,"message":"can't get planned budget"}`,
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
			expectedBody: `{"status":200,"body":{"actual_balance":100}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				expectedBudget := 100.0
				mockUsecase.EXPECT().GetCurrentBudget(gomock.Any()).Return(expectedBudget, nil)
			},
		},
		{
			name:         "Invalid userID",
			userID:       "invalidUserID",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"invalid url parameter"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				// No expectations for mockUsecase.
			},
		},
		{
			name:         "Internal server error",
			userID:       uuid.New().String(),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"status":500,"message":"can't get current budget"}`,
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
			expectedBody: `{"status":200,"body":{"account":[]}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().GetAccounts(gomock.Any()).Return([]models.Accounts{}, nil)
			},
		},
		{
			name:         "Invalid userID",
			userID:       "invalidUserID",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"invalid url parameter"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				// No expectations for mockUsecase.
			},
		},
		{
			name:         "No accounts found",
			userID:       uuidTest.String(),
			expectedCode: http.StatusOK,
			expectedBody: `{"status":200,"body":""}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				expectedError := models.NoSuchAccounts{}
				mockUsecase.EXPECT().GetAccounts(gomock.Any()).Return([]models.Accounts{}, &expectedError)
			},
		},
		{
			name:         "Internal server error",
			userID:       uuid.New().String(),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"status":500,"message":"no such account"}`,
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

func TestHandler_GetFeed(t *testing.T) {
	uuidTest := uuid.New()
	tests := []struct {
		name          string
		userID        string
		expectedCode  int
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
	}{
		{
			name:         "Successful call to GetFeed",
			userID:       uuid.New().String(),
			expectedCode: http.StatusOK,
			expectedBody: `{"status":200,"body":{"account":null,"balance":0,"planned_balance":0,"actual_balance":0,"err_message":""}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				userFeed := transfer_models.UserFeed{}
				mockUsecase.EXPECT().GetFeed(gomock.Any()).Return(userFeed, nil)
			},
		},
		{
			name:         "Invalid userID",
			userID:       "invalidUserID",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"invalid url parameter"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				// No expectations for mockUsecase.
			},
		},
		{
			name:         "No feed found",
			userID:       uuidTest.String(),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"no such feed info"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				expectedError := models.NoSuchAccounts{}
				userFeed := transfer_models.UserFeed{}
				mockUsecase.EXPECT().GetFeed(gomock.Any()).Return(userFeed, &expectedError)
			},
		},
		{
			name:         "Internal server error",
			userID:       uuid.New().String(),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"status":500,"message":"can't get feed info"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				internalError := errors.New("internal server error")
				userFeed := transfer_models.UserFeed{}
				mockUsecase.EXPECT().GetFeed(gomock.Any()).Return(userFeed, internalError)
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
			mockHandler.GetFeed(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_GetUser(t *testing.T) {
	uuidTest := uuid.New()
	tests := []struct {
		name          string
		userID        string
		expectedCode  int
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
	}{
		{
			name:         "Successful call to GetUser",
			userID:       uuid.New().String(),
			expectedCode: http.StatusOK,
			expectedBody: `{"status":200,"body":{"id":"00000000-0000-0000-0000-000000000000","username":"","planned_budget":0,"avatar_url":""}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {

				usr := &models.User{}
				mockUsecase.EXPECT().GetUser(gomock.Any()).Return(usr, nil)
			},
		},
		{
			name:         "Invalid userID",
			userID:       "invalidUserID",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"invalid url parameter"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				// No expectations for mockUsecase.
			},
		},
		{
			name:         "No user found",
			userID:       uuidTest.String(),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"no such user"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				expectedError := models.NoSuchUserError{}
				user := &models.User{}
				mockUsecase.EXPECT().GetUser(gomock.Any()).Return(user, &expectedError)
			},
		},
		{
			name:         "Internal server error",
			userID:       uuid.New().String(),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"status":500,"message":"can't get user"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				internalError := errors.New("internal server error")
				user := &models.User{}
				mockUsecase.EXPECT().GetUser(gomock.Any()).Return(user, internalError)
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
			mockHandler.Get(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}
