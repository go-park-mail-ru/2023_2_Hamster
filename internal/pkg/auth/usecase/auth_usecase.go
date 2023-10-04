package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secret = os.Getenv("SECRET")

type Usecase struct {
	authRepo auth.Repository
	userRepo user.Repository
	logger   logger.CustomLogger
}

type authClaims struct {
	UserID uuid.UUID `json:"id"`
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
func (u *Usecase) SignUpUser(user models.User) (uuid.UUID, auth.CookieToken, error) {
	salt := make([]byte, 8)
	rand.Read(salt)

	user.Salt = fmt.Sprintf("%x", salt)
	fmt.Println([]byte(user.Salt))
	user.Password = hashPassword(user.Password, salt)

	userId, err := u.userRepo.CreateUser(user)
	if err != nil {
		return uuid.Nil, auth.CookieToken{}, fmt.Errorf("[usecase] cannot create user: %w", err)
	}

	token, err := u.GenerateAccessToken(context.Background(), user)

	return userId, token, nil
}

func (u *Usecase) SignInUser(username, plainPassword string) (uuid.UUID, auth.CookieToken, error) {
	user, err := u.userRepo.GetUserByUsername(username)
	if err != nil {
		return uuid.Nil, auth.CookieToken{
			Value:   "",
			Expires: time.Now(),
		}, fmt.Errorf("[usecase] can't find user: %w", err)
	}

	salt, err := hex.DecodeString(user.Salt)
	if err != nil {
		return uuid.Nil, auth.CookieToken{
			Value:   "",
			Expires: time.Now(),
		}, fmt.Errorf("[usecase] salt from db decode error: %w", err)
	}

	hashedPassword := hashPassword(plainPassword, salt)
	if hashedPassword != user.Password {
		return uuid.Nil, auth.CookieToken{
			Value:   "",
			Expires: time.Now(),
		}, fmt.Errorf("[usecase] incorrect password")
	}

	token, err := u.GenerateAccessToken(context.Background(), *user)
	if err != nil {
		return uuid.Nil, auth.CookieToken{
			Value:   "",
			Expires: time.Now(),
		}, fmt.Errorf("[usecase] failed to generate access token: %w", err)
	}

	return user.ID, token, nil
}

// GetUserByCreds returns User if such exist in repository
func (u *Usecase) GetUserByCreds(ctx context.Context, username, plainPassword string) (*models.User, error) {
	user, err := u.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("[usecase] can't find user: %w", err)
	}

	return user, nil
}

// GetUserByAuthData returns User if such exist in repository
func (u *Usecase) GetUserByAuthData(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return u.userRepo.GetByID(userID)
}

func (u *Usecase) GenerateAccessToken(ctx context.Context, user models.User) (auth.CookieToken, error) {
	expTime := time.Now().UTC().Add(time.Hour * 24)

	tokenHeaderPayload := jwt.NewWithClaims(jwt.SigningMethodHS256, &authClaims{
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			Issuer:    "HammyWallet",
		},
	})

	tokenString, err := tokenHeaderPayload.SignedString([]byte(secret))
	if err != nil {
		return auth.CookieToken{
			Value:   "",
			Expires: time.Now(),
		}, err
	}

	return auth.CookieToken{
		Value:   tokenString,
		Expires: expTime,
	}, nil
}

func (u *Usecase) ValidateAccessToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return uuid.Nil, errors.New("[usecase] invalig signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("[usecase] invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["id"].(string)
		resID, err := uuid.Parse(userID)
		if err != nil {
			return uuid.Nil, fmt.Errorf("[usecase] invalid token: %w", err)
		}
		return resID, nil
	} else {
		return uuid.Nil, err
	}
}

// IncraseUserVersion inc User access token version
// func (u *Usecase) IncreaseUserVersion(ctx context.Context, userID uuid.UUID) error {
// 	if err := u.userRepo.IncreaseUserVersion(ctx, userID); err != nil {
// 		return fmt.Errorf("[usecase] error failed to update version: %w", err)
// 	}
// 	return nil
// }

// ChangePassword(ctx context.Context, userID uint32, password string) error

func hashPassword(pwd string, salt []byte) string {
	hash := sha256.New()
	hash.Write(append([]byte(pwd), salt...))
	return hex.EncodeToString(hash.Sum(nil))
}
