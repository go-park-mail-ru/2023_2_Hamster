package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetUserBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки для usecase и logger
	mockUsecase := mocks.NewMockUsecase(ctrl)
	//mockLogger := &logger.CustomLoggerMock{}

	// Создаем обработчик
	handler := &Handler{
		userService: mockUsecase,
		logger:      *logger.CreateCustomLogger(),
	}

	// Создаем UUID для тестов
	userID := uuid.New()

	// Сценарий 1: успешный вызов GetUserBalance
	expectedBalance := 100.0
	mockUsecase.EXPECT().GetUserBalance(userID).Return(expectedBalance, nil)
	url := "/api/user/" + userID.String() + "/balance"
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	handler.GetUserBalance(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, `{"Balance":100}`, recorder.Body.String())

	// Сценарий 2: ошибка валидации параметров
	invalidUserID := "invalidUserID" // неправильный формат UUID
	url = "api/user/" + invalidUserID + "/balance"
	req, err = http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	recorder = httptest.NewRecorder()
	handler.GetUserBalance(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Сценарий 3: ошибка GetUserBalance
	errorUserID := uuid.New()
	expectedError := models.NoSuchUserIdBalanceError{UserID: errorUserID}
	mockUsecase.EXPECT().GetUserBalance(errorUserID).Return(0.0, &expectedError)

	req, err = http.NewRequest("GET", "/user/balance?userID="+errorUserID.String(), nil)
	assert.NoError(t, err)

	recorder = httptest.NewRecorder()
	handler.GetUserBalance(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Сценарий 4: внутренняя ошибка сервера
	internalErrorUserID := uuid.New()
	internalError := errors.New("internal server error")
	mockUsecase.EXPECT().GetUserBalance(internalErrorUserID).Return(0.0, internalError)

	req, err = http.NewRequest("GET", "/user/balance?userID="+internalErrorUserID.String(), nil)
	assert.NoError(t, err)

	recorder = httptest.NewRecorder()
	handler.GetUserBalance(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}
