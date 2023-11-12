package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/hasher"
	logging "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

type Usecase struct {
	authRepo auth.Repository
	logger   logging.Logger
}

func NewUsecase(
	ar auth.Repository,
	log logging.Logger) *Usecase {
	return &Usecase{
		authRepo: ar,
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

	userId, err := u.authRepo.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("[usecase] cannot create user: %w", err)
	}

	user.ID = userId

	return userId, user.Username, nil
}

func (u *Usecase) Login(ctx context.Context, login, plainPassword string) (uuid.UUID, string, error) {
	user, err := u.authRepo.GetUserByLogin(ctx, login)
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

func (u *Usecase) CheckLoginUnique(ctx context.Context, login string) (bool, error) { // move from auth rep
	isUnique, err := u.authRepo.CheckLoginUnique(ctx, login)

	if err != nil {
		return false, fmt.Errorf("[usecase] can't login unique check")
	}

	return isUnique, nil
}

func (u *Usecase) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error) { // need test
	user, err := u.authRepo.GetByID(ctx, userID)
	if err != nil {
		return user, fmt.Errorf("[usecase] can't get user from repository %w", err)
	}
	return user, nil
}
