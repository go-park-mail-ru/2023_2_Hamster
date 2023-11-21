package http

// import (
// 	"context"
// 	"errors"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
// 	mocks "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/transaction/mocks"
// 	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
// 	"github.com/golang/mock/gomock"
// 	"github.com/google/uuid"
// 	"github.com/gorilla/mux"
// 	"github.com/stretchr/testify/assert"
// )

// func TestHandler_GetFeed(t *testing.T) {
// 	uuidTest := uuid.New()
// 	user := &models.User{ID: uuidTest}
// 	tests := []struct {
// 		name          string
// 		user          *models.User
// 		queryParam    string
// 		expectedCode  int
// 		expectedBody  string
// 		mockUsecaseFn func(*mocks.MockUsecase)
// 	}{
// 		{
// 			name:         "Successful call to GetFeed",
// 			user:         user,
// 			queryParam:   "page=2&page_size=10",
// 			expectedCode: http.StatusOK,
// 			expectedBody: `{"status":200,"body":{"transactions":[{"id":"00000000-0000-0000-0000-000000000000","account_income":"00000000-0000-0000-0000-000000000000","account_outcome":"00000000-0000-0000-0000-000000000000","income":0,"outcome":0,"date":"0001-01-01T00:00:00Z","payer":"","description":"","categories":null}],"is_all":true}}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 				mockUsecase.EXPECT().GetFeed(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Transaction{{UserID: uuidTest}}, true, nil)
// 			},
// 		},
// 		{
// 			name:         "Unauthorized Request",
// 			user:         nil,
// 			queryParam:   "page=2&page_size=10",
// 			expectedCode: http.StatusUnauthorized,
// 			expectedBody: `{"status":401,"message":"unauthorized"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 				// No service calls are expected for unauthorized request.
// 			},
// 		},
// 		{
// 			name:         "Invalid Query Parameter page_size null",
// 			user:         user,
// 			queryParam:   "page=10&page_size=''",
// 			expectedCode: http.StatusBadRequest,
// 			expectedBody: `{"status":400,"message":"invalid url parameter"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 			},
// 		},
// 		{
// 			name:         "Invalid Query Parameter page < 0",
// 			user:         user,
// 			queryParam:   "page=-1&page_size=''",
// 			expectedCode: http.StatusBadRequest,
// 			expectedBody: `{"status":400,"message":"invalid url parameter"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 			},
// 		},
// 		{
// 			name:         "Invalid Query Parameter page null",
// 			user:         user,
// 			queryParam:   "page=10&page_size=-1",
// 			expectedCode: http.StatusBadRequest,
// 			expectedBody: `{"status":400,"message":"invalid url parameter"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 			},
// 		},
// 		{
// 			name:         "Invalid Query Parameter page null",
// 			user:         user,
// 			queryParam:   "page='10'&page_size='0'",
// 			expectedCode: http.StatusBadRequest,
// 			expectedBody: `{"status":400,"message":"invalid url parameter"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 			},
// 		},
// 		{
// 			name:         "No Such Transaction Error",
// 			user:         user,
// 			queryParam:   "page=2&page_size=10",
// 			expectedCode: http.StatusOK,
// 			expectedBody: `{"status":204,"body":""}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 				errorNoSuchTransaction := models.NoSuchTransactionError{UserID: uuidTest}
// 				mockUsecase.EXPECT().GetFeed(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Transaction{}, false, &errorNoSuchTransaction)
// 			},
// 		},
// 		{
// 			name:         "Internal Server Error",
// 			user:         user,
// 			expectedCode: http.StatusInternalServerError,
// 			expectedBody: `{"status":500,"message":"can't get feed info"}`,
// 			mockUsecaseFn: func(mockService *mocks.MockUsecase) {
// 				internalServerError := errors.New("can't get feed info")
// 				mockService.EXPECT().GetFeed(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Transaction{}, true, internalServerError)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockService := mocks.NewMockUsecase(ctrl)
// 			tt.mockUsecaseFn(mockService)

// 			mockHandler := NewHandler(mockService, *logger.NewLogger(context.TODO()))

// 			req := httptest.NewRequest("GET", "/api/user/balance", nil)

// 			if tt.user != nil {
// 				ctx := context.WithValue(req.Context(), models.ContextKeyUserType{}, tt.user)
// 				req = req.WithContext(ctx)
// 				req.URL.RawQuery = tt.queryParam
// 			}

// 			recorder := httptest.NewRecorder()

// 			mockHandler.GetFeed(recorder, req)

// 			actual := strings.TrimSpace(recorder.Body.String())

// 			assert.Equal(t, tt.expectedCode, recorder.Code)
// 			assert.Equal(t, tt.expectedBody, actual)
// 		})
// 	}
// }

// func TestHandler_GetUserFromRequest(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		user          *models.User
// 		expectedCode  int
// 		expectedBody  string
// 		mockUsecaseFn func(*mocks.MockUsecase)
// 	}{
// 		{
// 			name:         "Unauthorized Request",
// 			user:         nil,
// 			expectedCode: http.StatusUnauthorized,
// 			expectedBody: `{"status":401,"message":"unauthorized"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockService := mocks.NewMockUsecase(ctrl)
// 			tt.mockUsecaseFn(mockService)

// 			mockHandler := NewHandler(mockService, *logger.NewLogger(context.TODO()))

// 			req := httptest.NewRequest("GET", "/api/user/balance", nil)

// 			ctx := context.WithValue(req.Context(), models.ContextKeyUserType{}, tt.user)
// 			req = req.WithContext(ctx)

// 			recorder := httptest.NewRecorder()

// 			mockHandler.GetFeed(recorder, req)

// 			actual := strings.TrimSpace(recorder.Body.String())

// 			assert.Equal(t, tt.expectedCode, recorder.Code)
// 			assert.Equal(t, tt.expectedBody, actual)
// 		})
// 	}
// }

// func TestHandler_CreateTransaction(t *testing.T) {
// 	uuidTest := uuid.New()
// 	user := &models.User{ID: uuidTest}

// 	tests := []struct {
// 		name          string
// 		user          *models.User
// 		requestBody   io.Reader
// 		expectedCode  int
// 		expectedBody  string
// 		mockUsecaseFn func(*mocks.MockUsecase)
// 	}{
// 		{
// 			name:         "Invalid JSON body",
// 			user:         user,
// 			requestBody:  strings.NewReader("invalid json data"),
// 			expectedCode: http.StatusBadRequest,
// 			expectedBody: `{"status":400,"message":"invalid input body"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 				// mockUsecase.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(uuidTest, nil)
// 			},
// 		},
// 		{
// 			name: "Successful Transaction Creation",
// 			user: user,
// 			requestBody: strings.NewReader(`{
// 			"account_income": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 			"account_outcome": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 			"categories": [],
// 			"date": "2023-10-02T15:30:00Z",
// 			"description": "string",
// 			"income": 0,
// 			"outcome": 0
// 		  }`),
// 			expectedCode: http.StatusOK,
// 			expectedBody: `{"status":200,"body":{"transaction_id":"` + uuidTest.String() + `"}}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 				mockUsecase.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(uuidTest, nil)
// 			},
// 		},
// 		{
// 			name: "Bad Check valid",
// 			user: user,
// 			requestBody: strings.NewReader(`{
// 				"account_income": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 				"account_outcome": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 				"categories": [],
// 				"date": "2023-10-02T15:30:00Z",
// 				"description": "string",
// 				"payer": "fffffffffffffffffffffffffffffffffffffffffffffffffff",
// 				"income": 0,
// 				"outcome": 0
// 			  }`),
// 			expectedCode:  http.StatusBadRequest,
// 			expectedBody:  `{"status":400,"message":"invalid input body"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {},
// 		},

// 		{
// 			name:         "Unauthorized Request",
// 			user:         nil,
// 			requestBody:  strings.NewReader(`{"field1": "value1", "field2": "value2"}`), // Replace with valid JSON
// 			expectedCode: http.StatusUnauthorized,
// 			expectedBody: `{"status":401,"message":"unauthorized"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 			},
// 		},
// 		{
// 			name: "Transaction Creation Error",
// 			user: user,
// 			requestBody: strings.NewReader(`{
// 				"account_income": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 				"account_outcome": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 				"categories": [],
// 				"date": "2023-10-02T15:30:00Z",
// 				"description": "string",
// 				"income": 0,
// 				"outcome": 0
// 			  }`), // Replace with valid JSON
// 			expectedCode: http.StatusBadRequest,
// 			expectedBody: `{"status":400,"message":"can't create transaction"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 				mockUsecase.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(uuidTest, errors.New("transaction not created"))
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockService := mocks.NewMockUsecase(ctrl)
// 			tt.mockUsecaseFn(mockService)

// 			mockHandler := NewHandler(mockService, *logger.NewLogger(context.TODO()))

// 			req := httptest.NewRequest("POST", "/api/transaction/create", tt.requestBody)

// 			if tt.user != nil {
// 				ctx := context.WithValue(req.Context(), models.ContextKeyUserType{}, tt.user)
// 				req = req.WithContext(ctx)
// 			}

// 			recorder := httptest.NewRecorder()

// 			mockHandler.Create(recorder, req)

// 			actual := strings.TrimSpace(recorder.Body.String())

// 			assert.Equal(t, tt.expectedCode, recorder.Code)
// 			assert.Equal(t, tt.expectedBody, actual)
// 		})
// 	}
// }

// func TestHandler_UpdateTransaction(t *testing.T) {
// 	uuidTest := uuid.New()
// 	user := &models.User{ID: uuidTest}

// 	tests := []struct {
// 		name          string
// 		user          *models.User
// 		requestBody   io.Reader
// 		expectedCode  int
// 		expectedBody  string
// 		mockUsecaseFn func(*mocks.MockUsecase)
// 	}{
// 		{
// 			name: "Successful Transaction Update",
// 			user: user,
// 			requestBody: strings.NewReader(`{
// 			"transaction_id": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 			"account_income": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 			"account_outcome": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 			"categories": [],
// 			"date": "2023-10-02T15:30:00Z",
// 			"description": "string",
// 			"income": 0,
// 			"outcome": 0
// 		}`),
// 			expectedCode: http.StatusOK,
// 			expectedBody: `{"status":200,"body":{}}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 				mockUsecase.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Return(nil)
// 			},
// 		},
// 		{
// 			name:         "Unauthorized Request",
// 			user:         nil,
// 			requestBody:  strings.NewReader(`{"field1": "value1", "field2": "value2"}`), // Replace with valid JSON
// 			expectedCode: http.StatusUnauthorized,
// 			expectedBody: `{"status":401,"message":"unauthorized"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 			},
// 		},
// 		{
// 			name:         "Invalid JSON body",
// 			user:         user,
// 			requestBody:  strings.NewReader("invalid json data"),
// 			expectedCode: http.StatusBadRequest,
// 			expectedBody: `{"status":400,"message":"invalid input body"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 				// mockUsecase.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Return(nil)
// 			},
// 		},
// 		{
// 			name: "Transaction Update Error",
// 			user: user,
// 			requestBody: strings.NewReader(`{
// 				"transaction_id": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 				"account_income": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 				"account_outcome": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 				"categories": [],
// 				"date": "2023-10-02T15:30:00Z",
// 				"description": "string",
// 				"income": 0,
// 				"outcome": 0
// 			}`),
// 			expectedCode: http.StatusInternalServerError,
// 			expectedBody: `{"status":500,"message":"can't get transaction"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 				mockUsecase.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Return(errors.New("transaction not updated"))
// 			},
// 		},
// 		{
// 			name: "Transaction Check valid",
// 			user: user,
// 			requestBody: strings.NewReader(`{
// 				"account_income": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 				"account_outcome": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 				"categories": [],
// 				"date": "2023-10-02T15:30:00Z",
// 				"description": "string",
// 				"payer": "fffffffffffffffffffffffffffffffffffffffffffffffffff",
// 				"income": 0,
// 				"outcome": 0
// 			  }`),
// 			expectedCode: http.StatusBadRequest,
// 			expectedBody: `{"status":400,"message":"invalid input body"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 			},
// 		},
// 		{
// 			name: "Transaction No Such User Error",
// 			user: user,
// 			requestBody: strings.NewReader(`{
// 				"transaction_id": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 				"account_income": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 				"account_outcome": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 				"categories": [],
// 				"date": "2023-10-02T15:30:00Z",
// 				"description": "string",
// 				"income": 0,
// 				"outcome": 0
// 			}`),
// 			expectedCode: http.StatusBadRequest,
// 			expectedBody: `{"status":400,"message":"can't such transactoin"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 				errorNoSuchTransaction := models.NoSuchTransactionError{UserID: uuidTest}
// 				mockUsecase.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Return(&errorNoSuchTransaction)
// 			},
// 		},
// 		{
// 			name: "User Forbidden",
// 			user: user,
// 			requestBody: strings.NewReader(`{
// 				"transaction_id": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 				"account_income": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 				"account_outcome": "7c62a6ef-2c4c-48c1-8a98-825fb6a3f0e6",
// 				"categories": [],
// 				"date": "2023-10-02T15:30:00Z",
// 				"description": "string",
// 				"income": 0,
// 				"outcome": 0
// 			}`),
// 			expectedCode: http.StatusForbidden,
// 			expectedBody: `{"status":403,"message":"user has no rights"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 				errorNoSuchUserForbidden := models.ForbiddenUserError{}
// 				mockUsecase.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Return(&errorNoSuchUserForbidden)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockService := mocks.NewMockUsecase(ctrl)
// 			tt.mockUsecaseFn(mockService)

// 			mockHandler := NewHandler(mockService, *logger.NewLogger(context.TODO()))

// 			req := httptest.NewRequest("POST", "/api/transaction/update", tt.requestBody)

// 			if tt.user != nil {
// 				ctx := context.WithValue(req.Context(), models.ContextKeyUserType{}, tt.user)
// 				req = req.WithContext(ctx)
// 			}

// 			recorder := httptest.NewRecorder()

// 			mockHandler.Update(recorder, req)

// 			actual := strings.TrimSpace(recorder.Body.String())

// 			assert.Equal(t, tt.expectedCode, recorder.Code)
// 			assert.Equal(t, tt.expectedBody, actual)
// 		})
// 	}
// }

// func TestHandler_TransactionDelete(t *testing.T) {
// 	uuidTest := uuid.New()
// 	user := &models.User{ID: uuidTest}
// 	tests := []struct {
// 		name          string
// 		userID        string
// 		expectedCode  int
// 		flag          bool
// 		expectedBody  string
// 		funcCtxUser   func(*models.User, context.Context) context.Context
// 		mockUsecaseFn func(*mocks.MockUsecase)
// 	}{
// 		{
// 			name:         "Successful call to GetUserBalance",
// 			userID:       uuid.New().String(),
// 			flag:         true,
// 			expectedCode: http.StatusOK,
// 			expectedBody: `{"status":200,"body":{}}`,

// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 				mockUsecase.EXPECT().DeleteTransaction(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
// 			},
// 			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
// 				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
// 			},
// 		},
// 		{
// 			name:   "Invalid userID",
// 			userID: "invalidUserID",
// 			flag:   true,

// 			expectedCode: http.StatusBadRequest,
// 			expectedBody: `{"status":400,"message":"invalid url parameter"}`,
// 			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
// 				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
// 			},
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 				// No expectations for mockUsecase.
// 			},
// 		},
// 		{
// 			name:         "Error No Such transaction error",
// 			userID:       uuidTest.String(),
// 			expectedCode: http.StatusBadRequest,
// 			flag:         true,
// 			expectedBody: `{"status":400,"message":"can't such transactoin"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 				expectedError := models.NoSuchTransactionError{}
// 				mockUsecase.EXPECT().DeleteTransaction(gomock.Any(), gomock.Any(), gomock.Any()).Return(&expectedError)
// 			},
// 			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
// 				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
// 			},
// 		},
// 		{
// 			name:         "Internal server error",
// 			userID:       uuid.New().String(),
// 			expectedCode: http.StatusInternalServerError,
// 			flag:         true,
// 			expectedBody: `{"status":500,"message":"cat't delete transaction"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 				//internalErrorUserID := uuid.New()
// 				internalError := errors.New("internal server error")
// 				mockUsecase.EXPECT().DeleteTransaction(gomock.Any(), gomock.Any(), gomock.Any()).Return(internalError)
// 			},
// 			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
// 				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
// 			},
// 		},
// 		{
// 			name:         "User Forbidden",
// 			userID:       uuid.New().String(),
// 			expectedCode: http.StatusForbidden,
// 			flag:         true,
// 			expectedBody: `{"status":403,"message":"user has no rights"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 				//internalErrorUserID := uuid.New()
// 				errorNoSuchUserForbidden := models.ForbiddenUserError{}
// 				mockUsecase.EXPECT().DeleteTransaction(gomock.Any(), gomock.Any(), gomock.Any()).Return(&errorNoSuchUserForbidden)
// 			},
// 			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
// 				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
// 			},
// 		},
// 		{
// 			name:         "Unauthorized",
// 			userID:       uuid.New().String(),
// 			expectedCode: http.StatusUnauthorized,
// 			flag:         false,
// 			expectedBody: `{"status":401,"message":"unauthorized"}`,
// 			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
// 			},
// 			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
// 				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockUsecase := mocks.NewMockUsecase(ctrl)
// 			tt.mockUsecaseFn(mockUsecase)

// 			mockHandler := NewHandler(mockUsecase, *logger.NewLogger(context.TODO()))

// 			url := "/api/transaction/" + tt.userID + "/delete"
// 			req := httptest.NewRequest("GET", url, nil)
// 			req = mux.SetURLVars(req, map[string]string{"transaction_id": tt.userID})

// 			if tt.flag {
// 				ctx := tt.funcCtxUser(user, req.Context())

// 				req = req.WithContext(ctx)
// 			}

// 			recorder := httptest.NewRecorder()
// 			mockHandler.Delete(recorder, req)

// 			actual := strings.TrimSpace(recorder.Body.String())

// 			assert.Equal(t, tt.expectedCode, recorder.Code)
// 			assert.Equal(t, tt.expectedBody, actual)
// 		})
// 	}
// }
