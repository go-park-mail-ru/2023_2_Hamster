package usecase_test

import (
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth"
	mock_auth "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth/usecase"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSignUpUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mock_auth.NewMockRepository(ctrl)
	mockUserRepo := mock_auth.NewMockRepository(ctrl)

	usecase := usecase.NewUsecase(mockAuthRepo, mockUserRepo, nil)

	user := models.User{
		// fill user details
	}

	mockUserRepo.EXPECT().CreateUser(gomock.Any()).Return(uuid.New(), nil)
	mockAuthRepo.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any()).Return(auth.CookieToken{}, nil)

	_, _, err := usecase.SignUpUser(user)
	assert.NoError(t, err)
}
