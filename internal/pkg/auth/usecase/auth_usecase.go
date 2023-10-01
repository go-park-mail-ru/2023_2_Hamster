package usecase

import (
	"context"
	"crypto/rand"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user"
	"github.com/golang-jwt/jwt/v5"
)

type Usecase struct {
	authRepo auth.Repository
	userRepo user.Repository
	logger   logger.CustomLogger
}

type authClaims struct {
	UserID      uint32 `json:"id"`
	UserVersion uint32 `json:"user_version"`
	jwt.RegisteredClaims
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
func (u *Usecase) SignUpUser(ctx context.Context, user models.User) (uint32, error) {
	salt := make([]byte, 8)

	rand.Read(salt)
}

// GetUserByCreds returns User if such exist in repository
func (u *Usecase) GetUserByCreds(ctx context.Context, username, plainPassword string) (*models.User, error) {

}

func (u *Usecase) LoginUser(username, plainPassword string) (string, error) {

}

// GetUserByAuthData returns User if such exist in repository
func (u *Usecase) GetUserByAuthData(ctx context.Context, userID, userVersion uint32) (*models.User, error) {

}

func (u *Usecase) GenerateAccessToken(ctx context.Context, userID, userVersion uint32) (string, error) {

}

func (u *Usecase) ValidateAccessToken(accessToken string) (uint32, uint32, error) {

}

// IncraseUserVersion inc User access token version
func (u *Usecase) IncreaseUserVersion(ctx context.Context, userID uint32) error {

}

// ChangePassword(ctx context.Context, userID uint32, password string) error
