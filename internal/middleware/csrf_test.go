package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	mock "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/csrf/mocks"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestCSRFSuccess(t *testing.T) {
	userID := uuid.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mock.NewMockUsecase(ctrl)

	uc.EXPECT().CheckCSRFToken(gomock.Any()).Return(userID, nil)

	handler := func(w http.ResponseWriter, r *http.Request) {}

	req := httptest.NewRequest(http.MethodPost, "http://www.your-domain.com", nil)

	req.Header.Set(csrfTokenHttpHeader, "2oo3hri3irj")
	user := &models.User{ID: userID}
	ctx := context.WithValue(req.Context(), models.ContextKeyUserType{}, user)
	req = req.WithContext(ctx)

	res := httptest.NewRecorder()
	handler(res, req)

	testM := NewCSRFMiddleware(uc, *logger.NewLogger(context.TODO()))
	middleware := testM.CheckCSRF(http.HandlerFunc(handler))

	middleware.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}

func TestCSRFUUIDNot(t *testing.T) {
	userID := uuid.New()
	userID1 := uuid.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mock.NewMockUsecase(ctrl)

	uc.EXPECT().CheckCSRFToken(gomock.Any()).Return(userID1, nil)

	handler := func(w http.ResponseWriter, r *http.Request) {}

	req := httptest.NewRequest(http.MethodPost, "http://www.your-domain.com", nil)

	req.Header.Set(csrfTokenHttpHeader, "2oo3hri3irj")
	user := &models.User{ID: userID}
	ctx := context.WithValue(req.Context(), models.ContextKeyUserType{}, user)
	req = req.WithContext(ctx)

	res := httptest.NewRecorder()
	handler(res, req)

	testM := NewCSRFMiddleware(uc, *logger.NewLogger(context.TODO()))
	middleware := testM.CheckCSRF(http.HandlerFunc(handler))

	middleware.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	actual := strings.TrimSpace(res.Body.String())
	expected := `{"status":400,"message":"invalid CSRF token"}`
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}

}

func TestCSRFCheckError(t *testing.T) {
	userID := uuid.New()
	userID1 := uuid.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mock.NewMockUsecase(ctrl)

	uc.EXPECT().CheckCSRFToken(gomock.Any()).Return(userID1, errors.New("err"))

	handler := func(w http.ResponseWriter, r *http.Request) {}

	req := httptest.NewRequest(http.MethodPost, "http://www.your-domain.com", nil)

	req.Header.Set(csrfTokenHttpHeader, "2oo3hri3irj")
	user := &models.User{ID: userID}
	ctx := context.WithValue(req.Context(), models.ContextKeyUserType{}, user)
	req = req.WithContext(ctx)

	res := httptest.NewRecorder()
	handler(res, req)

	testM := NewCSRFMiddleware(uc, *logger.NewLogger(context.TODO()))
	middleware := testM.CheckCSRF(http.HandlerFunc(handler))

	middleware.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	actual := strings.TrimSpace(res.Body.String())
	expected := `{"status":400,"message":"invalid CSRF token"}`
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}

}

func TestCSRFCheckBadToken(t *testing.T) {
	userID := uuid.New()
	//userID1 := uuid.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mock.NewMockUsecase(ctrl)

	//uc.EXPECT().CheckCSRFToken(gomock.Any()).Return(userID1, nil)

	handler := func(w http.ResponseWriter, r *http.Request) {}

	req := httptest.NewRequest(http.MethodPost, "http://www.your-domain.com", nil)

	req.Header.Set(csrfTokenHttpHeader, "")
	user := &models.User{ID: userID}
	ctx := context.WithValue(req.Context(), models.ContextKeyUserType{}, user)
	req = req.WithContext(ctx)

	res := httptest.NewRecorder()
	handler(res, req)

	testM := NewCSRFMiddleware(uc, *logger.NewLogger(context.TODO()))
	middleware := testM.CheckCSRF(http.HandlerFunc(handler))

	middleware.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	actual := strings.TrimSpace(res.Body.String())
	expected := `{"status":400,"message":"missing CSRF token"}`
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}

}

func TestCSRFUnauthorized(t *testing.T) {
	//userID1 := uuid.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mock.NewMockUsecase(ctrl)

	//uc.EXPECT().CheckCSRFToken(gomock.Any()).Return(userID1, nil)

	handler := func(w http.ResponseWriter, r *http.Request) {}

	req := httptest.NewRequest(http.MethodPost, "http://www.your-domain.com", nil)

	req.Header.Set(csrfTokenHttpHeader, "dd")

	res := httptest.NewRecorder()
	handler(res, req)

	testM := NewCSRFMiddleware(uc, *logger.NewLogger(context.TODO()))
	middleware := testM.CheckCSRF(http.HandlerFunc(handler))

	middleware.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	actual := strings.TrimSpace(res.Body.String())
	expected := `{"status":401,"message":"unauthorized"}`
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}

}
