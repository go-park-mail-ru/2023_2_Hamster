package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	proto "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/grpc/generated"
	mocks "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/mocks"
	mocksSession "github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions/mocks"
	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//func TestHandler_SignUp(t *testing.T) {
//	userid := uuid.New()
//	strUserId := userid.String()
//	tests := []struct {
//		name         string
//		requestBody  io.Reader
//		expectedCode int
//		expectedBody string
//		mockAU       func(*mocks.MockUsecase)
//		mockSU       func(*mocksSession.MockUsecase)
//	}{
//		{
//			name:         "Successful SignUp",
//			requestBody:  strings.NewReader(`{"login": "testlogin", "username": "testuser", "password": "testpassword"}`),
//			expectedCode: http.StatusOK,
//			expectedBody: fmt.Sprintf(`{"status":202,"body":{"id":"%s","username":"testuser"}}`, strUserId),
//			mockAU: func(mockAU *mocks.MockUsecase) {
//				mockAU.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(userid, "testuser", nil)
//			},
//			mockSU: func(mockSU *mocksSession.MockUsecase) {
//				mockSU.EXPECT().CreateSessionById(gomock.Any(), gomock.Any()).Return(models.Session{UserId: userid, Cookie: "testCookie"}, nil)
//			},
//		},
//		{
//			name:         "Corrupted request body",
//			requestBody:  strings.NewReader(`{"login": "testlogin","username": "testuser", "password": "testpassword`),
//			expectedCode: http.StatusBadRequest,
//			expectedBody: `{"status":400,"message":"Corrupted request body can't unmarshal"}`,
//			mockAU: func(mockAU *mocks.MockUsecase) {
//				// mockAU.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(userid, "testuser", nil)
//			},
//			mockSU: func(mockSU *mocksSession.MockUsecase) {
//				// mockSU.EXPECT().CreateSessionById(gomock.Any(), gomock.Any()).Return(models.Session{UserId: userid, Cookie: "testCookie"}, nil)
//			},
//		},
//		{
//			name:         "Error during SignUp",
//			requestBody:  strings.NewReader(`{"login": "testlogin","username": "testuser", "password": "testpassword"}`),
//			expectedCode: http.StatusTooManyRequests,
//			expectedBody: `{"status":429,"message":"Can't Sign Up user"}`,
//			mockAU: func(mockAU *mocks.MockUsecase) {
//				mockAU.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(uuid.Nil, "", errors.New("signup error"))
//			},
//			mockSU: func(mockSU *mocksSession.MockUsecase) {},
//		},
//		// Add more test cases as needed
//	}

func TestHandler_SignUp(t *testing.T) {
	userid := uuid.New()
	strUserId := userid.String()
	tests := []struct {
		name         string
		requestBody  io.Reader
		expectedCode int
		expectedBody string
		mockAU       func(client *mocks.MockAuthServiceClient)
		mockSU       func(*mocksSession.MockUsecase)
	}{
		{
			name:         "Successful SignUp",
			requestBody:  strings.NewReader(`{"login": "testlogin", "username": "testuser", "password": "testpassword"}`),
			expectedCode: http.StatusCreated,
			expectedBody: fmt.Sprintf(`{"status":201,"body":{"id":"%s","login":"testlogin","username":"testuser"}}`, strUserId),
			mockAU: func(mockAU *mocks.MockAuthServiceClient) {
				body := proto.SignUpResponseBody{
					Id:       userid.String(),
					Login:    "testlogin",
					Username: "testuser",
				}
				mockAU.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(&proto.SignUpResponse{
					Status: "200",
					Body:   &body,
				}, nil)
			},
			mockSU: func(mockSU *mocksSession.MockUsecase) {
				mockSU.EXPECT().CreateSessionById(gomock.Any(), gomock.Any()).Return(models.Session{UserId: userid, Cookie: "testCookie"}, nil)
			},
		},
		{
			name:         "Corrupted request body",
			requestBody:  strings.NewReader(`{"login": "testlogin","username": "testuser", "password": "testpassword`),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"Corrupted request body can't unmarshal"}`,
			mockAU: func(mockAU *mocks.MockAuthServiceClient) {
				// You may add more specific expectations if needed
				mockAU.EXPECT().SignUp(gomock.Any(), gomock.Any()).Times(0)
			},
			mockSU: func(mockSU *mocksSession.MockUsecase) {
				// You may add more specific expectations if needed
				mockSU.EXPECT().CreateSessionById(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:         "Error during SignUp",
			requestBody:  strings.NewReader(`{"login": "testlogin","username": "testuser", "password": "testpassword"}`),
			expectedCode: http.StatusTooManyRequests,
			expectedBody: `{"status":429,"message":"Can't Sign Up user"}`,
			mockAU: func(mockAU *mocks.MockAuthServiceClient) {
				mockAU.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil, errors.New("signup error"))
			},
			mockSU: func(mockSU *mocksSession.MockUsecase) {
				// You may add more specific expectations if needed
				mockSU.EXPECT().CreateSessionById(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAU := mocks.NewMockAuthServiceClient(ctrl)
			mockSU := mocksSession.NewMockUsecase(ctrl)
			tt.mockAU(mockAU)
			tt.mockSU(mockSU)

			handler := &Handler{
				client: mockAU,
				su:     mockSU,
				log:    *logger.NewLogger(context.TODO()),
			}

			req := httptest.NewRequest("POST", "/api/signup", tt.requestBody)
			recorder := httptest.NewRecorder()

			handler.SignUp(recorder, req)

			assert.Equal(t, tt.expectedCode, recorder.Result().StatusCode)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(recorder.Body.String()))
		})
	}
}

func TestHandler_Login(t *testing.T) {
	userID := uuid.New()
	strUserID := userID.String()
	tests := []struct {
		name         string
		requestBody  io.Reader
		expectedCode int
		expectedBody string
		mockAU       func(client *mocks.MockAuthServiceClient)
		mockSU       func(*mocksSession.MockUsecase)
	}{
		{
			name:         "Successful Login",
			requestBody:  strings.NewReader(`{"login": "testuser", "password": "testpassword"}`),
			expectedCode: http.StatusAccepted,
			expectedBody: fmt.Sprintf(`{"status":202,"body":{"id":"%s","login":"login","username":"testuser"}}`, strUserID),
			mockAU: func(mockAU *mocks.MockAuthServiceClient) {
				body := proto.LoginResponseBody{
					Id:       strUserID,
					Login:    "login",
					Username: "testuser",
				}
				mockAU.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Return(&proto.LoginResponse{
					Status: "200",
					Body:   &body,
				}, nil)
			},
			mockSU: func(mockSU *mocksSession.MockUsecase) {
				mockSU.EXPECT().CreateSessionById(gomock.Any(), gomock.Any()).Return(models.Session{UserId: userID, Cookie: "testCookie"}, nil)
			},
		},
		{
			name:         "Corrupted request body",
			requestBody:  strings.NewReader(`{"login": "testuser", "password": "testpassword`),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"message":"Corrupted request body can't unmarshal"}`,
			mockAU: func(mockAU *mocks.MockAuthServiceClient) {
				// mockAU.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Return(userID, "testuser", nil)
			},
			mockSU: func(mockSU *mocksSession.MockUsecase) {
				// mockSU.EXPECT().CreateSessionById(gomock.Any(), gomock.Any()).Return(models.Session{UserID: userID, Cookie: "testCookie"}, nil)
			},
		},
		{
			name:         "Error during Login",
			requestBody:  strings.NewReader(`{"login": "testuser", "password": "testpassword"}`),
			expectedCode: http.StatusTooManyRequests,
			expectedBody: `{"status":429,"message":"Can't Login user"}`,
			mockAU: func(mockAU *mocks.MockAuthServiceClient) {
				mockAU.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("login error"))
			},
			mockSU: func(mockSU *mocksSession.MockUsecase) {},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAU := mocks.NewMockAuthServiceClient(ctrl)
			mockSU := mocksSession.NewMockUsecase(ctrl)
			tt.mockAU(mockAU)
			tt.mockSU(mockSU)

			handler := &Handler{
				client: mockAU,
				su:     mockSU,
				log:    *logger.NewLogger(context.TODO()),
			}

			req := httptest.NewRequest("POST", "/api/login", tt.requestBody)
			recorder := httptest.NewRecorder()

			handler.Login(recorder, req)

			assert.Equal(t, tt.expectedCode, recorder.Result().StatusCode)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(recorder.Body.String()))
		})
	}
}

func TestHandler_HealthCheck(t *testing.T) {
	userID := uuid.New()
	strUserID := userID.String()
	sessionCookie := "testCookie"

	tests := []struct {
		name          string
		requestCookie *http.Cookie
		expectedCode  int
		expectedBody  string
		mockSU        func(*mocksSession.MockUsecase)
		mockAU        func(*mocks.MockAuthServiceClient)
	}{
		{
			name:          "Successful Health Check",
			requestCookie: &http.Cookie{Name: "session_id", Value: sessionCookie},
			expectedCode:  http.StatusOK,
			expectedBody:  fmt.Sprintf(`{"status":200,"body":{"id":"%s","login":"testlogin","username":"testuser"}}`, strUserID),
			mockSU: func(mockSU *mocksSession.MockUsecase) {
				mockSU.EXPECT().GetSessionByCookie(gomock.Any(), sessionCookie).Return(models.Session{UserId: userID, Cookie: sessionCookie}, nil)
			},
			mockAU: func(mockAU *mocks.MockAuthServiceClient) {
				body := proto.UserResponseBody{
					Id:       userID.String(),
					Login:    "testlogin",
					Username: "testuser",
				}
				mockAU.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&proto.UserResponse{
					Status: "200",
					Body:   &body,
				}, nil)
			},
		},
		{
			name:          "No Cookie Provided",
			requestCookie: nil,
			expectedCode:  http.StatusForbidden,
			expectedBody:  `{"status":403,"message":"No cookie provided"}`,
			mockSU:        func(mockSU *mocksSession.MockUsecase) {},
			mockAU:        func(mockAU *mocks.MockAuthServiceClient) {},
		},
		{
			name:          "Session Doesn't Exist",
			requestCookie: &http.Cookie{Name: "session_id", Value: sessionCookie},
			expectedCode:  http.StatusUnauthorized,
			expectedBody:  `{"status":401,"message":"Session doesn't exist login"}`,
			mockSU: func(mockSU *mocksSession.MockUsecase) {
				mockSU.EXPECT().GetSessionByCookie(gomock.Any(), sessionCookie).Return(models.Session{}, errors.New("session not found"))
			},
			mockAU: func(mockAU *mocks.MockAuthServiceClient) {},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSU := mocksSession.NewMockUsecase(ctrl)
			mockAU := mocks.NewMockAuthServiceClient(ctrl)
			tt.mockSU(mockSU)
			tt.mockAU(mockAU)

			handler := &Handler{
				client: mockAU,
				su:     mockSU,
				log:    *logger.NewLogger(context.TODO()),
			}

			req := httptest.NewRequest("POST", "/api/checkAuth", nil)
			if tt.requestCookie != nil {
				req.AddCookie(tt.requestCookie)
			}

			recorder := httptest.NewRecorder()

			handler.HealthCheck(recorder, req)

			assert.Equal(t, tt.expectedCode, recorder.Result().StatusCode)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(recorder.Body.String()))
		})
	}
}

func TestHandler_LogOut(t *testing.T) {
	sessionCookie := "testCookie"

	tests := []struct {
		name          string
		requestCookie *http.Cookie
		expectedCode  int
		expectedBody  string
		mockSU        func(*mocksSession.MockUsecase)
	}{
		{
			name:          "Successful Log Out",
			requestCookie: &http.Cookie{Name: "session_id", Value: sessionCookie},
			expectedCode:  http.StatusOK,
			expectedBody:  `{"status":200,"body":{}}`,
			mockSU: func(mockSU *mocksSession.MockUsecase) {
				mockSU.EXPECT().DeleteSessionByCookie(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:          "No Cookie Provided",
			requestCookie: nil,
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"status":400,"message":"No cookie provided"}`,
			mockSU:        func(mockSU *mocksSession.MockUsecase) {},
		},
		{
			name:          "Error Deleting Session",
			requestCookie: &http.Cookie{Name: "session_id", Value: sessionCookie},
			expectedCode:  http.StatusInternalServerError,
			expectedBody:  `{"status":500,"message":"Can't delete session"}`,
			mockSU: func(mockSU *mocksSession.MockUsecase) {
				mockSU.EXPECT().DeleteSessionByCookie(gomock.Any(), gomock.Any()).Return(errors.New("delete session error"))
			},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSU := mocksSession.NewMockUsecase(ctrl)
			tt.mockSU(mockSU)

			handler := &Handler{
				client: nil,
				su:     mockSU,
				log:    *logger.NewLogger(context.TODO()),
			}

			req := httptest.NewRequest("POST", "/api/logout", nil)
			if tt.requestCookie != nil {
				req.AddCookie(tt.requestCookie)
			}

			recorder := httptest.NewRecorder()

			handler.LogOut(recorder, req)

			assert.Equal(t, tt.expectedCode, recorder.Result().StatusCode)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(recorder.Body.String()))
		})
	}
}

func TestHandler_CheckLoginUnique(t *testing.T) {
	tests := []struct {
		name         string
		expectBody   io.Reader
		isUnique     bool
		expectedCode int
		expectedBody string
		mockAU       func(*mocks.MockAuthServiceClient)
	}{
		{
			name:         "Login is Unique",
			expectBody:   strings.NewReader(`{"login":"login"}`),
			isUnique:     true,
			expectedCode: http.StatusOK,
			expectedBody: `{"status":200,"body":{"status":"200","body":true}}`,
			mockAU: func(mockAU *mocks.MockAuthServiceClient) {
				mockAU.EXPECT().CheckLoginUnique(gomock.Any(), gomock.Any()).Return(&proto.UniqCheckResponse{
					Status: "200",
					Body:   true,
				}, nil)
			},
		},
		{
			name:         "Login is Not Unique",
			expectBody:   strings.NewReader(`{"login":"login"}`),
			isUnique:     false,
			expectedCode: http.StatusOK,
			expectedBody: `{"status":200,"body":{"status":"200"}}`,
			mockAU: func(mockAU *mocks.MockAuthServiceClient) {
				mockAU.EXPECT().CheckLoginUnique(gomock.Any(), gomock.Any()).Return(&proto.UniqCheckResponse{
					Status: "200",
					Body:   false,
				}, nil)
			},
		},
		{
			name:         "Error Checking Unique Login",
			expectBody:   strings.NewReader(`{"login":"login"}`),
			isUnique:     false,
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"status":500,"message":"Can't query DB"}`,
			mockAU: func(mockAU *mocks.MockAuthServiceClient) {
				mockAU.EXPECT().CheckLoginUnique(gomock.Any(), gomock.Any()).Return(&proto.UniqCheckResponse{
					Status: "200",
					Body:   false,
				}, errors.New("check unique error"))
			},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAU := mocks.NewMockAuthServiceClient(ctrl)
			tt.mockAU(mockAU)

			handler := &Handler{
				client: mockAU,
				su:     nil,
				log:    *logger.NewLogger(context.TODO()),
			}

			req := httptest.NewRequest("POST", "/api/loginCheck", tt.expectBody)
			recorder := httptest.NewRecorder()

			handler.CheckLoginUnique(recorder, req)

			assert.Equal(t, tt.expectedCode, recorder.Result().StatusCode)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(recorder.Body.String()))
		})
	}
}
