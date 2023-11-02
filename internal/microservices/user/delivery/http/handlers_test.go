package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/delivery/http/transfer_models"
	mocks "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/mocks"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/gorilla/mux"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetUserBalance(t *testing.T) {
	uuidTest := uuid.New()
	user := &models.User{ID: uuidTest}
	tests := []struct {
		name          string
		userID        string
		expectedCode  int
		funcCtxUser   func(*models.User, context.Context) context.Context
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
	}{
		{
			name:         "Successful call to GetUserBalance",
			userID:       uuid.New().String(),
			expectedCode: http.StatusOK,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedBody: `{"status":200,"body":{"balance":100}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				expectedBalance := 100.0
				mockUsecase.EXPECT().GetUserBalance(gomock.Any(), gomock.Any()).Return(expectedBalance, nil)
			},
		},
		{
			name:         "Error from GetUserBalance",
			userID:       uuidTest.String(),
			expectedCode: http.StatusBadRequest,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedBody: `{"status":400,"message":"no such balance"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				errorUserID := uuidTest
				expectedError := models.NoSuchUserIdBalanceError{UserID: errorUserID}
				mockUsecase.EXPECT().GetUserBalance(gomock.Any(), gomock.Any()).Return(0.0, &expectedError)
			},
		},
		{
			name:         "Internal server error",
			userID:       uuid.New().String(),
			expectedCode: http.StatusInternalServerError,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedBody: `{"status":500,"message":"can't get balance"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				//internalErrorUserID := uuid.New()
				internalError := errors.New("internal server error")
				mockUsecase.EXPECT().GetUserBalance(gomock.Any(), gomock.Any()).Return(0.0, internalError)
			},
		},
		{
			name:   "Unauthorized Request",
			userID: uuid.New().String(),

			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.Background()
			},
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {

			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"status":401,"message":"unauthorized"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			mockHandler := NewHandler(mockUsecase, logger.GetLogger())

			url := "/api/user/" + user.ID.String() + "/balance"
			req := httptest.NewRequest("GET", url, nil)
			ctx := tt.funcCtxUser(user, req.Context())

			req = req.WithContext(ctx)
			recorder := httptest.NewRecorder()
			req = req.WithContext(ctx)

			mockHandler.GetUserBalance(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_GetPlannedBudget(t *testing.T) {
	uuidTest := uuid.New()
	user := &models.User{ID: uuidTest}
	tests := []struct {
		name          string
		userID        string
		expectedCode  int
		funcCtxUser   func(*models.User, context.Context) context.Context
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
	}{
		{
			name:         "Successful call to GetPlannedBudget",
			userID:       uuid.New().String(),
			expectedCode: http.StatusOK,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedBody: `{"status":200,"body":{"planned_budget":500}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				expectedPlannedBudget := 500.0
				mockUsecase.EXPECT().GetPlannedBudget(gomock.Any(), gomock.Any()).Return(expectedPlannedBudget, nil)
			},
		},
		{
			name:   "Unauthorized Request",
			userID: uuid.New().String(),

			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.Background()
			},
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {

			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"status":401,"message":"unauthorized"}`,
		},
		{
			name:         "Error from GetPlannedBudget",
			userID:       uuidTest.String(),
			expectedCode: http.StatusBadRequest,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedBody: `{"status":400,"message":"no such planned budget"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				errorUserID := uuidTest
				expectedError := models.NoSuchPlannedBudgetError{UserID: errorUserID}
				mockUsecase.EXPECT().GetPlannedBudget(gomock.Any(), gomock.Any()).Return(0.0, &expectedError)
			},
		},
		{
			name:         "Internal server error",
			userID:       uuid.New().String(),
			expectedCode: http.StatusInternalServerError,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedBody: `{"status":500,"message":"can't get planned budget"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				internalError := errors.New("internal server error")
				mockUsecase.EXPECT().GetPlannedBudget(gomock.Any(), gomock.Any()).Return(0.0, internalError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			mockHandler := NewHandler(mockUsecase, logger.GetLogger())

			url := "/api/user/" + user.ID.String() + "/planned-budget"
			req := httptest.NewRequest("GET", url, nil)
			ctx := tt.funcCtxUser(user, req.Context())

			req = req.WithContext(ctx)
			recorder := httptest.NewRecorder()
			req = req.WithContext(ctx)

			mockHandler.GetPlannedBudget(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_GetCurrentBudget(t *testing.T) {
	uuidTest := uuid.New()
	user := &models.User{ID: uuidTest}
	tests := []struct {
		name          string
		userID        string
		expectedCode  int
		funcCtxUser   func(*models.User, context.Context) context.Context
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
	}{
		{
			name:         "Successful call to GetCurrentBudget",
			userID:       uuid.New().String(),
			expectedCode: http.StatusOK,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedBody: `{"status":200,"body":{"actual_budget":100}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				expectedBudget := 100.0
				mockUsecase.EXPECT().GetCurrentBudget(gomock.Any(), gomock.Any()).Return(expectedBudget, nil)
			},
		},
		{
			name:         "Unauthorized Request",
			userID:       "invalidUserID",
			expectedCode: http.StatusUnauthorized,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.Background()
			},
			expectedBody: `{"status":401,"message":"unauthorized"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				// No expectations for mockUsecase.
			},
		},
		{
			name:         "Internal server error",
			userID:       uuid.New().String(),
			expectedCode: http.StatusInternalServerError,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedBody: `{"status":500,"message":"can't get current budget"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				internalError := errors.New("internal server error")
				mockUsecase.EXPECT().GetCurrentBudget(gomock.Any(), gomock.Any()).Return(0.0, internalError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			mockHandler := NewHandler(mockUsecase, logger.GetLogger())

			url := "/api/user/" + tt.userID + "/current_budget"
			req := httptest.NewRequest("GET", url, nil)
			ctx := tt.funcCtxUser(user, req.Context())

			req = req.WithContext(ctx)
			recorder := httptest.NewRecorder()
			req = req.WithContext(ctx)

			mockHandler.GetCurrentBudget(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_GetAccounts(t *testing.T) {
	uuidTest := uuid.New()
	user := &models.User{ID: uuidTest}
	tests := []struct {
		name          string
		userID        string
		expectedCode  int
		funcCtxUser   func(*models.User, context.Context) context.Context
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
	}{
		{
			name:         "Successful call to GetAccounts",
			userID:       uuid.New().String(),
			expectedCode: http.StatusOK,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedBody: `{"status":200,"body":{"accounts":[]}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).Return([]models.Accounts{}, nil)
			},
		},
		{
			name:         "Unauthorized Request",
			userID:       "invalidUserID",
			expectedCode: http.StatusUnauthorized,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.Background()
			},
			expectedBody: `{"status":401,"message":"unauthorized"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				// No expectations for mockUsecase.
			},
		},
		{
			name:         "No accounts found",
			userID:       uuidTest.String(),
			expectedCode: http.StatusOK,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedBody: `{"status":204,"body":""}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				expectedError := models.NoSuchAccounts{}
				mockUsecase.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).Return([]models.Accounts{}, &expectedError)
			},
		},
		{
			name:         "Internal server error",
			userID:       uuid.New().String(),
			expectedCode: http.StatusInternalServerError,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedBody: `{"status":500,"message":"no such account"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				internalError := errors.New("internal server error")
				mockUsecase.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).Return([]models.Accounts{}, internalError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			mockHandler := NewHandler(mockUsecase, logger.GetLogger())

			url := "/api/user/" + tt.userID + "/accounts"
			req := httptest.NewRequest("GET", url, nil)
			ctx := tt.funcCtxUser(user, req.Context())

			req = req.WithContext(ctx)
			recorder := httptest.NewRecorder()
			req = req.WithContext(ctx)

			mockHandler.GetAccounts(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestHandler_GetFeed(t *testing.T) {
	uuidTest := uuid.New()
	user := &models.User{ID: uuidTest}
	tests := []struct {
		name          string
		userID        string
		expectedCode  int
		funcCtxUser   func(*models.User, context.Context) context.Context
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
	}{
		{
			name:         "Successful call to GetFeed",
			userID:       uuid.New().String(),
			expectedCode: http.StatusOK,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedBody: `{"status":200,"body":{"accounts":null,"balance":0,"planned_budget":0,"actual_budget":0}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				userFeed := &transfer_models.UserFeed{}
				mockUsecase.EXPECT().GetFeed(gomock.Any(), gomock.Any()).Return(userFeed, nil)
			},
		},
		{
			name:         "Unauthorized Request",
			userID:       "invalidUserID",
			expectedCode: http.StatusUnauthorized,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.Background()
			},
			expectedBody: `{"status":401,"message":"unauthorized"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				// No expectations for mockUsecase.
			},
		},
		{
			name:         "No feed found",
			userID:       uuidTest.String(),
			expectedCode: http.StatusBadRequest,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedBody: `{"status":400,"message":"no such feed info"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				expectedError := models.NoSuchAccounts{}
				userFeed := &transfer_models.UserFeed{}
				mockUsecase.EXPECT().GetFeed(gomock.Any(), gomock.Any()).Return(userFeed, &expectedError)
			},
		},
		{
			name:         "Internal server error",
			userID:       uuid.New().String(),
			expectedCode: http.StatusInternalServerError,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedBody: `{"status":500,"message":"can't get feed info"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				internalError := errors.New("internal server error")
				userFeed := &transfer_models.UserFeed{}
				mockUsecase.EXPECT().GetFeed(gomock.Any(), gomock.Any()).Return(userFeed, internalError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			mockHandler := NewHandler(mockUsecase, logger.GetLogger())

			url := "/api/user/" + tt.userID + "/accounts"
			req := httptest.NewRequest("GET", url, nil)
			ctx := tt.funcCtxUser(user, req.Context())

			req = req.WithContext(ctx)
			recorder := httptest.NewRecorder()
			req = req.WithContext(ctx)

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
			expectedBody: `{"status":200,"body":{"id":"00000000-0000-0000-0000-000000000000","login":"","username":"","planned_budget":0,"avatar_url":"00000000-0000-0000-0000-000000000000"}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {

				usr := &models.User{}
				mockUsecase.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(usr, nil)
			},
		},
		{
			name:         "Unauthorized Request",
			userID:       "invalidUserID",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"unauthorized"}`,
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
				mockUsecase.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(user, &expectedError)
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
				mockUsecase.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(user, internalError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			mockHandler := NewHandler(mockUsecase, logger.GetLogger())

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

func TestHandler_Update(t *testing.T) {
	testData := `{"username":"newUser", "planned_budget":123}`
	tests := []struct {
		name          string
		requestBody   io.Reader
		mockUsecaseFn func(*mocks.MockUsecase)
		user          *models.User
		funcCtxUser   func(*models.User, context.Context) context.Context
		expectedCode  int
		expectedBody  string
	}{
		{
			name:        "Successful Update",
			requestBody: strings.NewReader(testData),
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(nil)
			},
			user: &models.User{},
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"status":200,"body":{}}`,
		},
		{
			name:        "Invalid Request Body",
			requestBody: strings.NewReader("invalid json data"),
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				// No expectations for mockUsecase.
			},
			user: &models.User{Username: "user"},
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"invalid input body"}`,
		},
		{
			name:        "User Not Found",
			requestBody: strings.NewReader(testData),
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(&models.NoSuchUserError{})
			},
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			user:         &models.User{Username: "user"},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"no such user"}`,
		},
		{
			name:          "User Not Valid",
			requestBody:   strings.NewReader(`{"username":"newUserrrrrrrrrrrrrrrrrrrrrrrrrrrr", "planned_budget":111}`),
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {},
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			user:         &models.User{Username: "user"},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"invalid input body"}`,
		},
		{
			name:        "Internal Server Error",
			requestBody: strings.NewReader(testData),
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(errors.New("internal server error"))
			},
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			user:         &models.User{Username: "user"},
			expectedCode: http.StatusInternalServerError,

			expectedBody: `{"status":500,"message":"can't get user"}`,
		},
		{
			name:        "Unauthorized Request",
			requestBody: strings.NewReader(testData),
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
			},
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.Background()
			},
			user:         &models.User{},
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"status":401,"message":"unauthorized"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			mockHandler := NewHandler(mockUsecase, logger.GetLogger())

			url := "/api/user/update"
			req := httptest.NewRequest("PUT", url, tt.requestBody)

			ctx := tt.funcCtxUser(tt.user, req.Context())
			req = req.WithContext(ctx)
			recorder := httptest.NewRecorder()
			mockHandler.Update(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestUserHandler_UpdateProfilePhoto(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	//pathName := uuid.New().String()
	uuidTest := uuid.New()
	user := &models.User{ID: uuidTest}
	//usecaseMock := mocks.NewMockUsecase(ctrl)
	mockUsecase := mocks.NewMockUsecase(ctrl)
	tests := []struct {
		name           string
		userID         string
		funcCtxUser    func(*models.User, context.Context) context.Context
		mock           func() *http.Request
		expectedStatus int
	}{
		{
			name:   "OK",
			userID: uuid.New().String(),
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("path")
				_, err := partPath.Write([]byte(uuid.Nil.String()))
				if err != nil {
					t.Error(err)
				}

				part, err := writer.CreateFormFile("upload", "img.png")
				if err != nil {
					t.Error(err)
				}

				width := 200
				height := 100
				upLeft := image.Point{0, 0}
				lowRight := image.Point{width, height}

				img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

				// Colors are defined by Red, Green, Blue, Alpha uint8 values.
				cyan := color.RGBA{100, 200, 200, 0xff}

				// Set color for each pixel.
				for x := 0; x < width; x++ {
					for y := 0; y < height; y++ {
						switch {
						case x < width/2 && y < height/2: // upper left quadrant
							img.Set(x, y, cyan)
						case x >= width/2 && y >= height/2: // lower right quadrant
							img.Set(x, y, color.White)
						default:
							// Use zero value.
						}
					}
				}

				// Encode() takes an io.Writer. We pass the multipart field 'upload' that we defined
				// earlier which, in turn, writes to our io.Pipe
				err = png.Encode(part, img)
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					body)

				r.Header.Add("Content-Type", writer.FormDataContentType())

				mockUsecase.EXPECT().UpdatePhoto(gomock.Any(), gomock.Any()).Return(uuid.New(), nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Unauthorized Request",
			userID: "invaild_id",
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.Background()
			},
			mock: func() *http.Request {

				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					nil)
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "Error parse path",
			userID: uuid.New().String(),
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("path")
				_, err := partPath.Write([]byte("111"))
				if err != nil {
					t.Error(err)
				}

				part, err := writer.CreateFormFile("upload", "img.png")
				if err != nil {
					t.Error(err)
				}

				width := 200
				height := 100
				upLeft := image.Point{0, 0}
				lowRight := image.Point{width, height}

				img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

				// Colors are defined by Red, Green, Blue, Alpha uint8 values.
				cyan := color.RGBA{100, 200, 200, 0xff}

				// Set color for each pixel.
				for x := 0; x < width; x++ {
					for y := 0; y < height; y++ {
						switch {
						case x < width/2 && y < height/2: // upper left quadrant
							img.Set(x, y, cyan)
						case x >= width/2 && y >= height/2: // lower right quadrant
							img.Set(x, y, color.White)
						default:
							// Use zero value.
						}
					}
				}

				// Encode() takes an io.Writer. We pass the multipart field 'upload' that we defined
				// earlier which, in turn, writes to our io.Pipe
				err = png.Encode(part, img)
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					body)

				r.Header.Add("Content-Type", writer.FormDataContentType())

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Check large file",
			userID: uuid.New().String(),
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "No upload file",
			userID: uuid.New().String(),
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)
				go func() {
					defer writer.Close()

					partPath, _ := writer.CreateFormField("path")
					_, err := partPath.Write([]byte(uuid.Nil.String()))
					if err != nil {
						t.Error(err)
					}

				}()

				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					body)

				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Wrong data type",
			userID: uuid.New().String(),
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("path")
				_, err := partPath.Write([]byte("111"))
				if err != nil {
					t.Error(err)
				}

				part, err := writer.CreateFormFile("upload", "text.txt")
				if err != nil {
					t.Error(err)
				}

				_, err = part.Write([]byte("This is a text file content."))
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/user/updateProfilePhoto", body)
				r.Header.Add("Content-Type", writer.FormDataContentType())

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "File upload error",
			userID: uuid.New().String(),
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			mock: func() *http.Request {
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				defer writer.Close()

				_, err := writer.CreateFormFile("uploafd", "img.png")
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/user/updateProfilePhoto", body)
				r.Header.Add("Content-Type", writer.FormDataContentType())

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Internal Server",
			userID: uuid.New().String(),
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("path")
				_, err := partPath.Write([]byte(uuid.Nil.String()))
				if err != nil {
					t.Error(err)
				}

				part, err := writer.CreateFormFile("upload", "img.png")
				if err != nil {
					t.Error(err)
				}

				width := 200
				height := 100
				upLeft := image.Point{0, 0}
				lowRight := image.Point{width, height}

				img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

				cyan := color.RGBA{100, 200, 200, 0xff}

				for x := 0; x < width; x++ {
					for y := 0; y < height; y++ {
						switch {
						case x < width/2 && y < height/2: // upper left quadrant
							img.Set(x, y, cyan)
						case x >= width/2 && y >= height/2: // lower right quadrant
							img.Set(x, y, color.White)
						default:
						}
					}
				}

				err = png.Encode(part, img)
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					body)

				r.Header.Add("Content-Type", writer.FormDataContentType())

				mockUsecase.EXPECT().UpdatePhoto(gomock.Any(), gomock.Any()).Return(uuid.New(), errors.New("some err"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "No such User",
			userID: uuid.New().String(),
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("path")
				_, err := partPath.Write([]byte(uuid.Nil.String()))
				if err != nil {
					t.Error(err)
				}

				part, err := writer.CreateFormFile("upload", "img.png")
				if err != nil {
					t.Error(err)
				}

				width := 200
				height := 100
				upLeft := image.Point{0, 0}
				lowRight := image.Point{width, height}

				img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

				// Colors are defined by Red, Green, Blue, Alpha uint8 values.
				cyan := color.RGBA{100, 200, 200, 0xff}

				// Set color for each pixel.
				for x := 0; x < width; x++ {
					for y := 0; y < height; y++ {
						switch {
						case x < width/2 && y < height/2: // upper left quadrant
							img.Set(x, y, cyan)
						case x >= width/2 && y >= height/2: // lower right quadrant
							img.Set(x, y, color.White)
						default:
							// Use zero value.
						}
					}
				}

				// Encode() takes an io.Writer. We pass the multipart field 'upload' that we defined
				// earlier which, in turn, writes to our io.Pipe
				err = png.Encode(part, img)
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					body)

				r.Header.Add("Content-Type", writer.FormDataContentType())

				mockUsecase.EXPECT().UpdatePhoto(gomock.Any(), gomock.Any()).Return(uuid.New(), &models.NoSuchUserError{})
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			mockHandler := NewHandler(mockUsecase, logger.GetLogger())

			w := httptest.NewRecorder()
			r := test.mock()
			ctx := test.funcCtxUser(user, r.Context())
			r = r.WithContext(ctx)

			mockHandler.UpdatePhoto(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s, response %s", test.name, test.expectedStatus, w.Code, test.name, w.Body.String()))
		})
	}
}
