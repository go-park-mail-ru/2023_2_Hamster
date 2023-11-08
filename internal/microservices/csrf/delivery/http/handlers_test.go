package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	mocks "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/csrf/mocks"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Delivery(t *testing.T) {
	uuidTest := uuid.New()
	user := &models.User{ID: uuidTest}
	tests := []struct {
		name          string
		expectedCode  int
		funcCtxUser   func(*models.User, context.Context) context.Context
		expectedBody  string
		mockUsecaseFn func(*mocks.MockUsecase)
	}{
		{
			name:         "Successful call to CSRF",
			expectedCode: http.StatusOK,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedBody: `{"status":200,"body":{"csrf":"token_valid"}}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().GenerateCSRFToken(gomock.Any()).Return("token_valid", nil)
			},
		},
		{
			name:         "Internal server error",
			expectedCode: http.StatusInternalServerError,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.WithValue(ctx, models.ContextKeyUserType{}, user)
			},
			expectedBody: `{"status":500,"message":"failed to get CSRF-token"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				mockUsecase.EXPECT().GenerateCSRFToken(gomock.Any()).Return("", errors.New("err"))
			},
		},
		{
			name:         "Unauthorized",
			expectedCode: http.StatusUnauthorized,
			funcCtxUser: func(user *models.User, ctx context.Context) context.Context {
				return context.Background()
			},
			expectedBody: `{"status":401,"message":"unauthorized"}`,
			mockUsecaseFn: func(mockUsecase *mocks.MockUsecase) {
				//internalErrorUserID := uuid.New()
				//internalError := errors.New("internal server error")
				//mockUsecase.EXPECT().GetUserBalance(gomock.Any(), gomock.Any()).Return(0.0, internalError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			mockHandler := NewHandler(mockUsecase, *logger.NewLogger(context.TODO()))

			url := "/api/user/" + user.ID.String() + "/balance"
			req := httptest.NewRequest("GET", url, nil)
			ctx := tt.funcCtxUser(user, req.Context())

			req = req.WithContext(ctx)
			recorder := httptest.NewRecorder()
			req = req.WithContext(ctx)

			mockHandler.GetCSRF(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}
