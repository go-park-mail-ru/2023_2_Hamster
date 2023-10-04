package usecase_test

import (
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	mock_auth "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth/usecase"
	mock_user "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user/mocks"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSignUpUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger.CreateCustomLogger()

	mockAuthRepo := mock_auth.NewMockRepository(ctrl)
	mockUserRepo := mock_user.NewMockRepository(ctrl)

	usecase := usecase.NewUsecase(mockAuthRepo, mockUserRepo, *mockLogger)

	user := models.User{
		Username: "user",
		Password: "12345",
	}

	mockUserRepo.EXPECT().CreateUser(gomock.Any()).Return(uuid.New(), nil)

	_, _, err := usecase.SignUpUser(user)
	assert.NoError(t, err)
}
