package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	genAccount "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account/delivery/grpc/generated"
	mocks "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/account/mocks"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateAccount(t *testing.T) {
	uuidTest := uuid.New()
	user := &models.User{ID: uuidTest}
	tests := []struct {
		name           string
		user           *models.User
		expectedCode   int
		expectedBody   string
		mockUsecaseFn  func(*mocks.MockAccountServiceClient)
		requestPayload string
	}{
		{
			name:         "Successful Account Creation",
			user:         user,
			expectedCode: http.StatusOK,
			expectedBody: `{"status":200,"body":{"account_id":"` + uuidTest.String() + `"}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockAccountServiceClient) {
				mockUsecase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&genAccount.CreateAccountResponse{AccountId: uuidTest.String()}, nil)
			},
			requestPayload: `{"balance": 100, "accumulation": true, "balance_enabled": true, "mean_payment": "monthly"}`,
		},
		{
			name:           "Unauthorized Request",
			user:           nil,
			expectedCode:   http.StatusUnauthorized,
			expectedBody:   `{"status":401,"message":"unauthorized"}`,
			mockUsecaseFn:  func(mockUsecase *mocks.MockAccountServiceClient) {},
			requestPayload: `{"balance": 100, "accumulation": true, "balance_enabled": true, "mean_payment": "monthly"}`,
		},
		{
			name:           "Invalid Request Body",
			user:           user,
			expectedCode:   http.StatusBadRequest,
			expectedBody:   `{"status":400,"message":"invalid input body"}`,
			mockUsecaseFn:  func(mockUsecase *mocks.MockAccountServiceClient) {},
			requestPayload: `invalid_json`,
		},
		{
			name:         "Error in Account Creation",
			user:         user,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"can't create account"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockAccountServiceClient) {
				mockUsecase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("account not created"))
			},
			requestPayload: `{"balance": 100, "accumulation": true, "balance_enabled": true, "mean_payment": "monthly"}`,
		},
		{
			name:         "Error Valid Check",
			user:         user,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"invalid input body"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockAccountServiceClient) {
			},
			requestPayload: `{"balance": 100, "accumulation": true, "balance_enabled": true, "meanPayment": "monthly"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockAccountServiceClient(ctrl)
			tt.mockUsecaseFn(mockService)

			mockHandler := NewHandler(mockService, *logger.NewLogger(context.TODO()))

			req := httptest.NewRequest("POST", "/api/account/create", strings.NewReader(tt.requestPayload))
			req.Header.Set("Content-Type", "application/json")

			if tt.user != nil {
				ctx := context.WithValue(req.Context(), models.ContextKeyUserType{}, tt.user)
				req = req.WithContext(ctx)
			}

			recorder := httptest.NewRecorder()

			mockHandler.Create(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_UpdateAccount(t *testing.T) {
	uuidTest := uuid.New()
	user := &models.User{ID: uuidTest}
	tests := []struct {
		name           string
		user           *models.User
		expectedCode   int
		expectedBody   string
		mockUsecaseFn  func(*mocks.MockAccountServiceClient)
		requestPayload string
	}{
		{
			name:         "Successful Account Update",
			user:         user,
			expectedCode: http.StatusOK,
			expectedBody: `{"status":200,"body":{}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockAccountServiceClient) {
				mockUsecase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, nil)
			},
			requestPayload: `{"id": "` + uuidTest.String() + `", "balance": 150, "accumulation": true, "balance_enabled": true, "mean_payment": "monthly"}`,
		},
		{
			name:           "Unauthorized Request",
			user:           nil,
			expectedCode:   http.StatusUnauthorized,
			expectedBody:   `{"status":401,"message":"unauthorized"}`,
			mockUsecaseFn:  func(mockUsecase *mocks.MockAccountServiceClient) {},
			requestPayload: `{"id": "` + uuidTest.String() + `", "balance": 150, "accumulation": true, "balanceEnabled": true, "meanPayment": "monthly"}`,
		},
		{
			name:           "Invalid Request Body",
			user:           user,
			expectedCode:   http.StatusBadRequest,
			expectedBody:   `{"status":400,"message":"invalid input body"}`,
			mockUsecaseFn:  func(mockUsecase *mocks.MockAccountServiceClient) {},
			requestPayload: `invalid_json`,
		},
		{
			name:         "Error in Account Update - NoSuchAccount",
			user:         user,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"can't such account"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockAccountServiceClient) {
				mockUsecase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, &models.NoSuchAccounts{})
			},
			requestPayload: `{"id": "` + uuidTest.String() + `", "balance": 150, "accumulation": true, "balanceEnabled": true, "meanPayment": "monthly"}`,
		},
		{
			name:         "Error in Account Update - ForbiddenUserError",
			user:         user,
			expectedCode: http.StatusForbidden,
			expectedBody: `{"status":403,"message":"user has no rights"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockAccountServiceClient) {
				mockUsecase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, &models.ForbiddenUserError{})
			},
			requestPayload: `{"id": "` + uuidTest.String() + `", "balance": 150, "accumulation": true, "balanceEnabled": true, "meanPayment": "monthly"}`,
		},
		{
			name:         "Error in Account Update - Some Errors",
			user:         user,
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"status":500,"message":"can't get account"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockAccountServiceClient) {
				mockUsecase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, errors.New("err"))
			},
			requestPayload: `{"id": "` + uuidTest.String() + `", "balance": 150, "accumulation": true, "balanceEnabled": true, "meanPayment": "monthly"}`,
		},
		{
			name:         "Error in Account Update - Check valid",
			user:         user,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"invalid input body"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockAccountServiceClient) {
			},
			requestPayload: `{"idd": "ff", "balance": 150, "accumulation": true, "balanceEnabled": true, "meaenPayment": "monthly"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockAccountServiceClient(ctrl)
			tt.mockUsecaseFn(mockService)

			mockHandler := NewHandler(mockService, *logger.NewLogger(context.TODO()))

			req := httptest.NewRequest("POST", "/api/account/update", strings.NewReader(tt.requestPayload))
			req.Header.Set("Content-Type", "application/json")

			if tt.user != nil {
				ctx := context.WithValue(req.Context(), models.ContextKeyUserType{}, tt.user)
				req = req.WithContext(ctx)
			}

			recorder := httptest.NewRecorder()

			mockHandler.Update(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_DeleteAccount(t *testing.T) {
	uuidTest := uuid.New()
	user := &models.User{ID: uuidTest}
	tests := []struct {
		name           string
		user           *models.User
		expectedCode   int
		expectedBody   string
		accountUrl     string
		mockUsecaseFn  func(*mocks.MockAccountServiceClient)
		requestPayload string
	}{
		{
			name:         "Successful Account Deletion",
			user:         user,
			expectedCode: http.StatusOK,
			expectedBody: `{"status":200,"body":{}}`,
			accountUrl:   "account_id",
			mockUsecaseFn: func(mockUsecase *mocks.MockAccountServiceClient) {
				mockUsecase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, nil)
			},
			requestPayload: `{"accountID": "` + uuidTest.String() + `"}`,
		},
		{
			name:           "Unauthorized Request",
			user:           nil,
			expectedCode:   http.StatusUnauthorized,
			expectedBody:   `{"status":401,"message":"unauthorized"}`,
			accountUrl:     "account_id",
			mockUsecaseFn:  func(mockUsecase *mocks.MockAccountServiceClient) {},
			requestPayload: `{"accountID": "` + uuidTest.String() + `"}`,
		},
		{
			name:         "Error in Account Deletion - NoSuchAccount",
			user:         user,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"can't such account"}`,
			accountUrl:   "account_id",
			mockUsecaseFn: func(mockUsecase *mocks.MockAccountServiceClient) {
				mockUsecase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, &models.NoSuchAccounts{})
			},
			requestPayload: `{"accountID": "` + uuidTest.String() + `"}`,
		},
		{
			name:         "Error in Account Deletion - ForbiddenUserError",
			user:         user,
			expectedCode: http.StatusForbidden,
			expectedBody: `{"status":403,"message":"user has no rights"}`,
			accountUrl:   "account_id",
			mockUsecaseFn: func(mockUsecase *mocks.MockAccountServiceClient) {
				mockUsecase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, &models.ForbiddenUserError{})
			},
			requestPayload: `{"accountID": "` + uuidTest.String() + `"}`,
		},
		{
			name:         "Error in Account Deletion - Some Errors",
			user:         user,
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"status":500,"message":"can't get account"}`,
			accountUrl:   "account_id",
			mockUsecaseFn: func(mockUsecase *mocks.MockAccountServiceClient) {
				mockUsecase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, errors.New("err"))
			},
			requestPayload: `{"accountID": "` + uuidTest.String() + `"}`,
		},
		{
			name:         "Error in Account Deletion - Invalid Url Parameter",
			user:         user,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"invalid url parameter"}`,
			accountUrl:   "account",
			mockUsecaseFn: func(mockUsecase *mocks.MockAccountServiceClient) {
			},
			requestPayload: `{"accountID": "` + uuidTest.String() + `"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockAccountServiceClient(ctrl)
			tt.mockUsecaseFn(mockService)

			mockHandler := NewHandler(mockService, *logger.NewLogger(context.TODO()))

			url := "/api/account/delete"
			req := httptest.NewRequest("GET", url, nil)
			req = mux.SetURLVars(req, map[string]string{tt.accountUrl: uuidTest.String()})

			if tt.user != nil {
				ctx := context.WithValue(req.Context(), models.ContextKeyUserType{}, tt.user)
				req = req.WithContext(ctx)
			}

			recorder := httptest.NewRecorder()
			mockHandler.Delete(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}
