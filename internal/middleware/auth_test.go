package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	mockR "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user/mocks"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	mockS "github.com/go-park-mail-ru/2023_2_Hamster/internal/monolithic/sessions/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestAuthenticationMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	su := mockS.NewMockUsecase(ctrl)
	ur := mockR.NewMockRepository(ctrl)
	logger := logger.NewLogger(context.TODO())

	sessionID := "valid_session_id"
	userID := uuid.New()
	user := &models.User{ID: userID}

	su.EXPECT().GetSessionByCookie(gomock.Any(), sessionID).Return(models.Session{UserId: userID}, nil)
	ur.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "http://www.your-domain.com", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})

	res := httptest.NewRecorder()

	testM := NewAuthMiddleware(su, ur, *logger)
	middleware := testM.Authentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	middleware.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("middleware returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestNilCookies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	su := mockS.NewMockUsecase(ctrl)
	ur := mockR.NewMockRepository(ctrl)
	logger := logger.NewLogger(context.TODO())

	sessionID := ""
	//userID := uuid.New()
	//user := &models.User{ID: userID}

	//su.EXPECT().GetSessionByCookie(gomock.Any(), sessionID).Return(models.Session{UserId: userID}, nil)
	//ur.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "http://www.your-domain.com", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})

	res := httptest.NewRecorder()

	testM := NewAuthMiddleware(su, ur, *logger)
	middleware := testM.Authentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	middleware.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusUnauthorized {
		t.Errorf("middleware returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	actual := strings.TrimSpace(res.Body.String())
	expected := `{"status":401,"message":"missing token unauthorized"}`
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}

}

func TestIvalidCookies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	su := mockS.NewMockUsecase(ctrl)
	ur := mockR.NewMockRepository(ctrl)
	logger := logger.NewLogger(context.TODO())

	//sessionID := ""
	//userID := uuid.New()
	//user := &models.User{ID: userID}

	//su.EXPECT().GetSessionByCookie(gomock.Any(), sessionID).Return(models.Session{UserId: userID}, nil)
	//ur.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "http://www.your-domain.com", nil)
	//req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})

	res := httptest.NewRecorder()

	testM := NewAuthMiddleware(su, ur, *logger)
	middleware := testM.Authentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	middleware.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusUnauthorized {
		t.Errorf("middleware returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	actual := strings.TrimSpace(res.Body.String())
	expected := `{"status":401,"message":"missing token unauthorized"}`
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}

}

func TestGetSessionByCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	su := mockS.NewMockUsecase(ctrl)
	ur := mockR.NewMockRepository(ctrl)
	logger := logger.NewLogger(context.TODO())

	sessionID := "valid_session_id"
	userID := uuid.New()
	//user := &models.User{ID: userID}

	su.EXPECT().GetSessionByCookie(gomock.Any(), sessionID).Return(models.Session{UserId: userID}, errors.New("err"))
	//ur.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "http://www.your-domain.com", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})

	res := httptest.NewRecorder()

	testM := NewAuthMiddleware(su, ur, *logger)
	middleware := testM.Authentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	middleware.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusUnauthorized {
		t.Errorf("middleware returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	actual := strings.TrimSpace(res.Body.String())
	expected := `{"status":401,"message":"token validation failed unauthorized"}`
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}
}

func TestGetSessionGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	su := mockS.NewMockUsecase(ctrl)
	ur := mockR.NewMockRepository(ctrl)
	logger := logger.NewLogger(context.TODO())

	sessionID := "valid_session_id"
	userID := uuid.New()
	//user := &models.User{ID: userID}

	su.EXPECT().GetSessionByCookie(gomock.Any(), gomock.Any()).Return(models.Session{UserId: userID}, nil)
	ur.EXPECT().GetByID(gomock.Any(), userID).Return(nil, errors.New("err"))

	req := httptest.NewRequest(http.MethodGet, "http://www.your-domain.com", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})

	res := httptest.NewRecorder()

	testM := NewAuthMiddleware(su, ur, *logger)
	middleware := testM.Authentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	middleware.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusUnauthorized {
		t.Errorf("middleware returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	actual := strings.TrimSpace(res.Body.String())
	expected := `{"status":401,"message":"userAuth check failed"}`
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}
}
