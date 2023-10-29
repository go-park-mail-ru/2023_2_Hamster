package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/hasher"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Usecase struct {
	authRepo auth.Repository
	userRepo user.Repository
	logger   logger.CustomLogger
}

func NewUsecase(
	ar auth.Repository,
	ur user.Repository,
	log logger.CustomLogger) *Usecase {
	return &Usecase{
		authRepo: ar,
		userRepo: ur,
		logger:   log,
	}
}

// SignUpUser creates new User and returns it's id
func (u *Usecase) SignUp(ctx context.Context, input auth.SignUpInput) (uuid.UUID, string, error) {
	var user models.User

	ok, err := u.authRepo.CheckLoginUnique(ctx, input.Login)
	if err != nil {
		u.logger.Error("Error checking login uniqueness: ", err)
		return uuid.Nil, "", fmt.Errorf("[usecase] error checking login uniqueness: %w", err)
	}

	if !ok {
		u.logger.Error("Login already exist ", input.Login)
		return uuid.Nil, "", fmt.Errorf("[usecase] username already exist")
	}

	hash, err := hasher.GeneratePasswordHash(input.PlaintPassword)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("[usecase] hash func error: %v", err)
	}

	user.Login = input.Login
	user.Username = input.Username
	user.Password = hash

	userId, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("[usecase] cannot create user: %w", err)
	}

	user.ID = userId

	return userId, user.Username, nil
}

func (u *Usecase) Login(ctx context.Context, login, plainPassword string) (uuid.UUID, string, error) {
	user, err := u.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("[usecase] can't find user: %w", err)
	}

	ok, err := hasher.VerfiyPassword(plainPassword, user.Password)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("[usecase] Password Comparation Error: %w", err)
	}
	if !ok {
		return uuid.Nil, "", fmt.Errorf("[usecase] incorrect password")
	}

	return user.ID, user.Username, nil
}

func (u *Usecase) IsLoginUnique(ctx context.Context, login string) (bool, error) { // TODO: move to auth repo
	isUnique, err := u.authRepo.CheckLoginUnique(ctx, login)

	if err != nil {
		return false, fmt.Errorf("[usecase] can`t login unique check")
	}

	return isUnique, nil
}
