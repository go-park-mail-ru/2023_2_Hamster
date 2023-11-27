package http

/*
func TestHandler_SignUp(t *testing.T) {
	userid := uuid.New()
	strUserId := userid.String()
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
			requestBody:  strings.NewReader(`{"login": "testlogin", "username": "testuser", "password": "testpassword"}`),
			expectedCode: http.StatusOK,
			expectedBody: fmt.Sprintf(`{"status":202,"body":{"id":"%s","username":"testuser"}}`, strUserId),
			mockAU: func(mockAU *mocks.MockUsecase) {
				mockAU.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(userid, "testuser", nil)
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
			mockAU: func(mockAU *mocks.MockUsecase) {
				// mockAU.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(userid, "testuser", nil)
			},
			mockSU: func(mockSU *mocksSession.MockUsecase) {
				// mockSU.EXPECT().CreateSessionById(gomock.Any(), gomock.Any()).Return(models.Session{UserId: userid, Cookie: "testCookie"}, nil)
			},
		},
		{
			name:         "Error during SignUp",
			requestBody:  strings.NewReader(`{"login": "testlogin","username": "testuser", "password": "testpassword"}`),
			expectedCode: http.StatusTooManyRequests,
			expectedBody: `{"status":429,"message":"Can't Sign Up user"}`,
			mockAU: func(mockAU *mocks.MockUsecase) {
				mockAU.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(uuid.Nil, "", errors.New("signup error"))
			},
			mockSU: func(mockSU *mocksSession.MockUsecase) {},
		},
		// Add more test cases as needed
	}
// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	mocks "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/mocks"
// 	mocksSession "github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions/mocks"
// 	"github.com/google/uuid"

// 	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
// 	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// )

// func TestHandler_SignUp(t *testing.T) {
// 	userid := uuid.New()
// 	strUserId := userid.String()
// 	tests := []struct {
// 		name         string
// 		requestBody  io.Reader
// 		expectedCode int
// 		expectedBody string
// 		mockAU       func(*mocks.MockUsecase)
// 		mockSU       func(*mocksSession.MockUsecase)
// 	}{
// 		{
// 			name:         "Successful SignUp",
// 			requestBody:  strings.NewReader(`{"login": "testlogin", "username": "testuser", "password": "testpassword"}`),
// 			expectedCode: http.StatusOK,
// 			expectedBody: fmt.Sprintf(`{"status":202,"body":{"id":"%s","username":"testuser"}}`, strUserId),
// 			mockAU: func(mockAU *mocks.MockUsecase) {
// 				mockAU.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(userid, "testuser", nil)
// 			},
// 			mockSU: func(mockSU *mocksSession.MockUsecase) {
// 				mockSU.EXPECT().CreateSessionById(gomock.Any(), gomock.Any()).Return(models.Session{UserId: userid, Cookie: "testCookie"}, nil)
// 			},
// 		},
// 		{
// 			name:         "Corrupted request body",
// 			requestBody:  strings.NewReader(`{"login": "testlogin","username": "testuser", "password": "testpassword`),
// 			expectedCode: http.StatusBadRequest,
// 			expectedBody: `{"status":400,"message":"Corrupted request body can't unmarshal"}`,
// 			mockAU: func(mockAU *mocks.MockUsecase) {
// 				// mockAU.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(userid, "testuser", nil)
// 			},
// 			mockSU: func(mockSU *mocksSession.MockUsecase) {
// 				// mockSU.EXPECT().CreateSessionById(gomock.Any(), gomock.Any()).Return(models.Session{UserId: userid, Cookie: "testCookie"}, nil)
// 			},
// 		},
// 		{
// 			name:         "Error during SignUp",
// 			requestBody:  strings.NewReader(`{"login": "testlogin","username": "testuser", "password": "testpassword"}`),
// 			expectedCode: http.StatusTooManyRequests,
// 			expectedBody: `{"status":429,"message":"Can't Sign Up user"}`,
// 			mockAU: func(mockAU *mocks.MockUsecase) {
// 				mockAU.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(uuid.Nil, "", errors.New("signup error"))
// 			},
// 			mockSU: func(mockSU *mocksSession.MockUsecase) {},
// 		},
// 		// Add more test cases as needed
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockAU := mocks.NewMockUsecase(ctrl)
// 			mockSU := mocksSession.NewMockUsecase(ctrl)
// 			tt.mockAU(mockAU)
// 			tt.mockSU(mockSU)

// 			handler := &Handler{
// 				au:  mockAU,
// 				su:  mockSU,
// 				log: *logger.NewLogger(context.TODO()),
// 			}

// 			req := httptest.NewRequest("POST", "/api/signup", tt.requestBody)
// 			recorder := httptest.NewRecorder()

// 			handler.SignUp(recorder, req)

// 			assert.Equal(t, tt.expectedCode, recorder.Result().StatusCode)
// 			assert.Equal(t, tt.expectedBody, strings.TrimSpace(recorder.Body.String()))
// 		})
// 	}
// }

// func TestHandler_Login(t *testing.T) {
// 	userID := uuid.New()
// 	strUserID := userID.String()
// 	tests := []struct {
// 		name         string
// 		requestBody  io.Reader
// 		expectedCode int
// 		expectedBody string
// 		mockAU       func(*mocks.MockUsecase)
// 		mockSU       func(*mocksSession.MockUsecase)
// 	}{
// 		{
// 			name:         "Successful Login",
// 			requestBody:  strings.NewReader(`{"login": "testuser", "password": "testpassword"}`),
// 			expectedCode: http.StatusOK,
// 			expectedBody: fmt.Sprintf(`{"status":202,"body":{"id":"%s","username":"testuser"}}`, strUserID),
// 			mockAU: func(mockAU *mocks.MockUsecase) {
// 				mockAU.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Return(userID, "testuser", nil)
// 			},
// 			mockSU: func(mockSU *mocksSession.MockUsecase) {
// 				mockSU.EXPECT().CreateSessionById(gomock.Any(), gomock.Any()).Return(models.Session{UserId: userID, Cookie: "testCookie"}, nil)
// 			},
// 		},
// 		{
// 			name:         "Corrupted request body",
// 			requestBody:  strings.NewReader(`{"login": "testuser", "password": "testpassword`),
// 			expectedCode: http.StatusBadRequest,
// 			expectedBody: `{"status":400,"message":"Corrupted request body can't unmarshal"}`,
// 			mockAU: func(mockAU *mocks.MockUsecase) {
// 				// mockAU.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Return(userID, "testuser", nil)
// 			},
// 			mockSU: func(mockSU *mocksSession.MockUsecase) {
// 				// mockSU.EXPECT().CreateSessionById(gomock.Any(), gomock.Any()).Return(models.Session{UserID: userID, Cookie: "testCookie"}, nil)
// 			},
// 		},
// 		{
// 			name:         "Error during Login",
// 			requestBody:  strings.NewReader(`{"login": "testuser", "password": "testpassword"}`),
// 			expectedCode: http.StatusTooManyRequests,
// 			expectedBody: `{"status":429,"message":"Can't Login user"}`,
// 			mockAU: func(mockAU *mocks.MockUsecase) {
// 				mockAU.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, "", errors.New("login error"))
// 			},
// 			mockSU: func(mockSU *mocksSession.MockUsecase) {},
// 		},
// 		// Add more test cases as needed
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockAU := mocks.NewMockUsecase(ctrl)
// 			mockSU := mocksSession.NewMockUsecase(ctrl)
// 			tt.mockAU(mockAU)
// 			tt.mockSU(mockSU)

// 			handler := &Handler{
// 				au:  mockAU,
// 				su:  mockSU,
// 				log: *logger.NewLogger(context.TODO()),
// 			}

// 			req := httptest.NewRequest("POST", "/api/login", tt.requestBody)
// 			recorder := httptest.NewRecorder()

// 			handler.Login(recorder, req)

// 			assert.Equal(t, tt.expectedCode, recorder.Result().StatusCode)
// 			assert.Equal(t, tt.expectedBody, strings.TrimSpace(recorder.Body.String()))
// 		})
// 	}
// }

// /*
// 	func TestHandler_HealthCheck(t *testing.T) {
// 		userID := uuid.New()
// 		strUserID := userID.String()
// 		sessionCookie := "testCookie"

		tests := []struct {
			name          string
			requestCookie *http.Cookie
			expectedCode  int
			expectedBody  string
			mockSU        func(*mocksSession.MockUsecase)
		}{
			{
				name:          "Successful Health Check",
				requestCookie: &http.Cookie{Name: "session_id", Value: sessionCookie},
				expectedCode:  http.StatusOK,
				expectedBody:  fmt.Sprintf(`{"status":200,"body":{"user_id":"%s","cookie":"%s"}}`, strUserID, sessionCookie),
				expectedBody:  fmt.Sprintf(`{"status":200,"body":{"id":"%s","username":"%s"}}`, strUserID, sessionCookie),
				mockSU: func(mockSU *mocksSession.MockUsecase) {
					mockSU.EXPECT().GetSessionByCookie(gomock.Any(), sessionCookie).Return(models.Session{UserId: userID, Cookie: sessionCookie}, nil)
				},
			},
			{
				name:          "No Cookie Provided",
				requestCookie: nil,
				expectedCode:  http.StatusForbidden,
				expectedBody:  `{"status":403,"message":"No cookie provided"}`,
				mockSU:        func(mockSU *mocksSession.MockUsecase) {},
			},
			{
				name:          "Session Doesn't Exist",
				requestCookie: &http.Cookie{Name: "session_id", Value: sessionCookie},
				expectedCode:  http.StatusUnauthorized,
				expectedBody:  `{"status":401,"message":"Session doesn't exist login"}`,
				mockSU: func(mockSU *mocksSession.MockUsecase) {
					mockSU.EXPECT().GetSessionByCookie(gomock.Any(), sessionCookie).Return(models.Session{}, errors.New("session not found"))
				},
			},
			// Add more test cases as needed
		}
// 		tests := []struct {
// 			name          string
// 			requestCookie *http.Cookie
// 			expectedCode  int
// 			expectedBody  string
// 			mockSU        func(*mocksSession.MockUsecase)
// 		}{
// 			{
// 				name:          "Successful Health Check",
// 				requestCookie: &http.Cookie{Name: "session_id", Value: sessionCookie},
// 				expectedCode:  http.StatusOK,
// <<<<<<< HEAD
// 				expectedBody:  fmt.Sprintf(`{"status":200,"body":{"user_id":"%s","cookie":"%s"}}`, strUserID, sessionCookie),
// =======
// 				expectedBody:  fmt.Sprintf(`{"status":200,"body":{"id":"%s","username":"%s"}}`, strUserID, sessionCookie),
// >>>>>>> 493d329e5b644d4e30dec179c1f48f05106223bb
// 				mockSU: func(mockSU *mocksSession.MockUsecase) {
// 					mockSU.EXPECT().GetSessionByCookie(gomock.Any(), sessionCookie).Return(models.Session{UserId: userID, Cookie: sessionCookie}, nil)
// 				},
// 			},
// 			{
// 				name:          "No Cookie Provided",
// 				requestCookie: nil,
// 				expectedCode:  http.StatusForbidden,
// 				expectedBody:  `{"status":403,"message":"No cookie provided"}`,
// 				mockSU:        func(mockSU *mocksSession.MockUsecase) {},
// 			},
// 			{
// 				name:          "Session Doesn't Exist",
// 				requestCookie: &http.Cookie{Name: "session_id", Value: sessionCookie},
// 				expectedCode:  http.StatusUnauthorized,
// 				expectedBody:  `{"status":401,"message":"Session doesn't exist login"}`,
// 				mockSU: func(mockSU *mocksSession.MockUsecase) {
// 					mockSU.EXPECT().GetSessionByCookie(gomock.Any(), sessionCookie).Return(models.Session{}, errors.New("session not found"))
// 				},
// 			},
// 			// Add more test cases as needed
// 		}

// 		for _, tt := range tests {
// 			t.Run(tt.name, func(t *testing.T) {
// 				ctrl := gomock.NewController(t)
// 				defer ctrl.Finish()

// 				mockSU := mocksSession.NewMockUsecase(ctrl)
// 				tt.mockSU(mockSU)

// 				handler := &Handler{
// 					au:  nil,
// 					su:  mockSU,
// 					log: *logger.NewLogger(context.TODO()),
// 				}

// 				req := httptest.NewRequest("GET", "/api/health", nil)
// 				if tt.requestCookie != nil {
// 					req.AddCookie(tt.requestCookie)
// 				}

// 				recorder := httptest.NewRecorder()

// 				handler.HealthCheck(recorder, req)

				assert.Equal(t, tt.expectedCode, recorder.Result().StatusCode)
				assert.Equal(t, tt.expectedBody, strings.TrimSpace(recorder.Body.String()))
			})
		}
	}
*/
/*
func TestHandler_LogOut(t *testing.T) {
	sessionCookie := "testCookie"
// 				assert.Equal(t, tt.expectedCode, recorder.Result().StatusCode)
// 				assert.Equal(t, tt.expectedBody, strings.TrimSpace(recorder.Body.String()))
// 			})
// 		}
// 	}
// */
// func TestHandler_LogOut(t *testing.T) {
// 	sessionCookie := "testCookie"

// 	tests := []struct {
// 		name          string
// 		requestCookie *http.Cookie
// 		expectedCode  int
// 		expectedBody  string
// 		mockSU        func(*mocksSession.MockUsecase)
// 	}{
// 		{
// 			name:          "Successful Log Out",
// 			requestCookie: &http.Cookie{Name: "session_id", Value: sessionCookie},
// 			expectedCode:  http.StatusOK,
// 			expectedBody:  `{"status":200,"body":{}}`,
// 			mockSU: func(mockSU *mocksSession.MockUsecase) {
// 				mockSU.EXPECT().DeleteSessionByCookie(gomock.Any(), gomock.Any()).Return(nil)
// 			},
// 		},
// 		{
// 			name:          "No Cookie Provided",
// 			requestCookie: nil,
// 			expectedCode:  http.StatusBadRequest,
// 			expectedBody:  `{"status":400,"message":"No cookie provided"}`,
// 			mockSU:        func(mockSU *mocksSession.MockUsecase) {},
// 		},
// 		{
// 			name:          "Error Deleting Session",
// 			requestCookie: &http.Cookie{Name: "session_id", Value: sessionCookie},
// 			expectedCode:  http.StatusInternalServerError,
// 			expectedBody:  `{"status":500,"message":"Can't delete session"}`,
// 			mockSU: func(mockSU *mocksSession.MockUsecase) {
// 				mockSU.EXPECT().DeleteSessionByCookie(gomock.Any(), gomock.Any()).Return(errors.New("delete session error"))
// 			},
// 		},
// 		// Add more test cases as needed
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockSU := mocksSession.NewMockUsecase(ctrl)
// 			tt.mockSU(mockSU)

// 			handler := &Handler{
// 				au:  nil,
// 				su:  mockSU,
// 				log: *logger.NewLogger(context.TODO()),
// 			}

// 			req := httptest.NewRequest("POST", "/api/logout", nil)
// 			if tt.requestCookie != nil {
// 				req.AddCookie(tt.requestCookie)
// 			}

// 			recorder := httptest.NewRecorder()

// 			handler.LogOut(recorder, req)

// 			assert.Equal(t, tt.expectedCode, recorder.Result().StatusCode)
// 			assert.Equal(t, tt.expectedBody, strings.TrimSpace(recorder.Body.String()))
// 		})
// 	}
// }

// func TestHandler_CheckLoginUnique(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		urlParam     string
// 		isUnique     bool
// 		expectedCode int
// 		expectedBody string
// 		mockAU       func(*mocks.MockUsecase)
// 	}{
// 		{
// 			name:         "Login is Unique",
// 			urlParam:     "uniqueLogin",
// 			isUnique:     true,
// 			expectedCode: http.StatusOK,
// 			expectedBody: `{"status":200,"body":true}`,
// 			mockAU: func(mockAU *mocks.MockUsecase) {
// 				mockAU.EXPECT().CheckLoginUnique(gomock.Any(), gomock.Any()).Return(true, nil)
// 			},
// 		},
// 		{
// 			name:         "Login is Not Unique",
// 			urlParam:     "nonUniqueLogin",
// 			isUnique:     false,
// 			expectedCode: http.StatusOK,
// 			expectedBody: `{"status":200,"body":false}`,
// 			mockAU: func(mockAU *mocks.MockUsecase) {
// 				mockAU.EXPECT().CheckLoginUnique(gomock.Any(), gomock.Any()).Return(false, nil)
// 			},
// 		},
// 		{
// 			name:         "Error Checking Unique Login",
// 			urlParam:     "errorLogin",
// 			isUnique:     false,
// 			expectedCode: http.StatusInternalServerError,
// 			expectedBody: `{"status":500,"message":"Can't get unique info login"}`,
// 			mockAU: func(mockAU *mocks.MockUsecase) {
// 				mockAU.EXPECT().CheckLoginUnique(gomock.Any(), gomock.Any()).Return(false, errors.New("check unique error"))
// 			},
// 		},
// 		// Add more test cases as needed
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockAU := mocks.NewMockUsecase(ctrl)
// 			tt.mockAU(mockAU)

// 			handler := &Handler{
// 				au:  mockAU,
// 				su:  nil,
// 				log: *logger.NewLogger(context.TODO()),
// 			}

// 			req := httptest.NewRequest("GET", fmt.Sprintf("/api/check-login-unique?%s=%s", userloginUrlParam, tt.urlParam), nil)
// 			recorder := httptest.NewRecorder()

// 			handler.CheckLoginUnique(recorder, req)
/*
			assert.Equal(t, tt.expectedCode, recorder.Result().StatusCode)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(recorder.Body.String()))
		})
	}
}
*/
// 			assert.Equal(t, tt.expectedCode, recorder.Result().StatusCode)
// 			assert.Equal(t, tt.expectedBody, strings.TrimSpace(recorder.Body.String()))
// 		})
// 	}
// }
