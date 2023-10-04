package usecase_test

import (
	"testing"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth/usecase"
)

func TestSignUpUser(t *testing.T) {
	user := models.User{
		Username: "user",
		Password: "123456",
	}

	ar := New

	au := usecase.NewUsecase()
}
