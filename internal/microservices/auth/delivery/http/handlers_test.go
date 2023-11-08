package http

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mocks "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/mocks"
	mocksSession "github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions/mocks"
	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_SignUp(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  io.Reader
		expectedCode int
		expectedBody string
		mockAU       func(*mocks.MockUsecase)
		mockSU       func(*mocksSession.MockUsecase)
	}{
		{
			name:         "Successful SignUp",
			requestBody:  strings.NewReader(`{"username": "testuser", "password": "testpassword"}`),
			expectedCode: http.StatusAccepted,
			expectedBody: `{"id": "testUserID", "username": "testuser"}`,
			mockAU: func(mockAU *mocks.MockUsecase) {
				mockAU.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(uuid.New(), "testuser", nil)
			},
			mockSU: func(mockSU *mocksSession.MockUsecase) {
				mockSU.EXPECT().CreateSessionById(gomock.Any(), gomock.Any()).Return(models.Session{UserId: uuid.New(), Cookie: "testCookie"}, nil)
			},
		},
		{
			name:         "Error during SignUp",
			requestBody:  strings.NewReader(`{"username": "testuser", "password": "testpassword"}`),
			expectedCode: http.StatusTooManyRequests,
			expectedBody: `{"status":429,"message":"Can't Sign Up user"}`,
			mockAU: func(mockAU *mocks.MockUsecase) {
				mockAU.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(uuid.Nil, "", errors.New("signup error"))
			},
			mockSU: func(mockSU *mocksSession.MockUsecase) {},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAU := mocks.NewMockUsecase(ctrl)
			mockSU := mocksSession.NewMockUsecase(ctrl)
			tt.mockAU(mockAU)
			tt.mockSU(mockSU)

			handler := &Handler{
				au:  mockAU,
				su:  mockSU,
				log: *logger.NewLogger(context.TODO()),
			}

			req := httptest.NewRequest("POST", "/api/signup", tt.requestBody)
			recorder := httptest.NewRecorder()

			handler.SignUp(recorder, req)

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(recorder.Body.String()))
		})
	}
}
