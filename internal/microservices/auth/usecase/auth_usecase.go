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
	logger   logging.Logger // legacy
}

func NewUsecase(
	ar auth.Repository,
	log logging.Logger) *Usecase {
	return &Usecase{
		authRepo: ar,
		logger:   log, // legacy
	}
}

// SignUpUser creates new User and returns it's id
func (u *Usecase) SignUp(ctx context.Context, input auth.SignUpInput) (uuid.UUID, string, string, error) {
	var user models.User

	ok, err := u.authRepo.CheckLoginUnique(ctx, input.Login)
	if err != nil { // Db error
		u.logger.Error("Error checking login uniqueness: ", err)
		return uuid.Nil, "", "", fmt.Errorf("[usecase] error checking login uniqueness: %w", err)
	}

	if !ok {
		u.logger.Error("Login already exist ", input.Login)
		return uuid.Nil, "", "", fmt.Errorf("[usecase] %w", &models.UserAlreadyExistsError{}) // Error login exist
	}

	hash, err := hasher.GeneratePasswordHash(input.PlaintPassword)
	if err != nil {
		return uuid.Nil, "", "", fmt.Errorf("[usecase] hash func error: %w", err) // Hash error
	}

	user.Login = input.Login
	user.Username = input.Username
	user.Password = hash

	userId, err := u.authRepo.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, "", "", fmt.Errorf("[usecase] cannot create user: %w", err)
	}

	user.ID = userId

	return userId, user.Login, user.Username, nil
}

func (u *Usecase) Login(ctx context.Context, login, plainPassword string) (uuid.UUID, string, string, error) {
	user, err := u.authRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return uuid.Nil, "", "", fmt.Errorf("[usecase] can't find user: %w", err)
	}

	ok, err := hasher.VerfiyPassword(plainPassword, user.Password)
	if err != nil {
		return uuid.Nil, "", "", fmt.Errorf("[usecase] Password Comparation Error: %w", err)
	}
	if !ok {
		return uuid.Nil, "", "", fmt.Errorf("[usecase] password hash doesn't match the real one: %w", &models.IncorrectPasswordError{UserID: user.ID})
	}

	return user.ID, user.Login, user.Username, nil
}

func (u *Usecase) CheckLoginUnique(ctx context.Context, login string) (bool, error) {
	isUnique, err := u.authRepo.CheckLoginUnique(ctx, login)
	if err != nil {
		return false, fmt.Errorf("[usecase] can't login unique check")
	}

	return isUnique, nil
}

func (u *Usecase) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := u.authRepo.GetByID(ctx, userID)
	if err != nil {
		return user, fmt.Errorf("[usecase] can't get user from repository %w", err)
	}
	return user, nil
}

func (u *Usecase) ChangePassword(ctx context.Context, input auth.ChangePasswordInput) error {
	user, err := u.authRepo.GetByID(ctx, input.Id)
	if err != nil {
		return fmt.Errorf("[usecase] can't find user: %w", err)
	}

	ok, err := hasher.VerfiyPassword(input.OldPassword, user.Password)
	if err != nil {
		return fmt.Errorf("[usecase] Password Comparation Error: %w", err)
	}
	if !ok {
		return fmt.Errorf("[usecase] password hash doesn't match the real one: %w", &models.IncorrectPasswordError{UserID: user.ID})
	}

	newpwd, err := hasher.GeneratePasswordHash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("[usecase] can't change password intenal err: %w", err)
	}

	if err := u.authRepo.ChangePassword(ctx, input.Id, newpwd); err != nil {
		return fmt.Errorf("[usecase] can't change password %w", err)
	}
	return nil
}
