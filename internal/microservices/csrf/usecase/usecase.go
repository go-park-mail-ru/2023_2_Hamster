package usecase

import (
	"errors"
	"fmt"
	"os"
	"time"

	logging "github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var tokenSecret = os.Getenv("SECRET")

const csrfTokenTTL = time.Minute //30 * time.Minute
const jwtParsingMaxTime = 3 * time.Second

type Usecase struct {
	logger logging.Logger
}

func NewUsecase(
	log logging.Logger) *Usecase {
	return &Usecase{
		logger: log,
	}
}

type jwtCSRFClaims struct {
	UserId uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

func (u *Usecase) GenerateCSRFToken(userID uuid.UUID) (string, error) {
	claims := &jwtCSRFClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(csrfTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	signedToken, err := signToken(claims, tokenSecret)
	if err != nil {
		return "", fmt.Errorf("[usecase] failed to sign acess token: %w", err)
	}

	return signedToken, nil
}

func signToken(claims jwt.Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func checkToken(tokenStr string, secret string, claims jwt.Claims) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("[usecase] invalid token signing method")
			}
			return []byte(secret), nil
		}, jwt.WithLeeway(jwtParsingMaxTime))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (u *Usecase) CheckCSRFToken(acessToken string) (uuid.UUID, error) {
	token, err := checkToken(acessToken, tokenSecret, &jwtCSRFClaims{})
	if err != nil {
		return uuid.Nil, fmt.Errorf("[usecase] invalid CSRF token")
	}

	claims, ok := token.Claims.(*jwtCSRFClaims)
	if !ok {
		return uuid.Nil, errors.New("[usecase] token claims are not of type *tokenClaims")
	}

	now := time.Now().UTC()
	if claims.ExpiresAt.Time.Before(now) {
		return uuid.Nil, errors.New("[usecase] csrf token is expired")
	}

	return claims.UserId, nil
}
