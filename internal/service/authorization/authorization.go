package authorization_service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	initdata "github.com/telegram-mini-apps/init-data-golang"
	"github.com/yazmeyaa/ton_game/internal/config"
)

type AuthorizationService interface {
	ValidateTelegramData(data string) error
	ValidateJWT(token string) error
	ParseJWTClaims(token string) (jwt.MapClaims, error)
	GenerateToken(userID int) (string, error)
}

type authorizationService struct {
	appConfig *config.Config
}

// ParseJWTClaims implements AuthorizationService.
func (a *authorizationService) ParseJWTClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(a.appConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// GenerateToken implements AuthorizationService.
func (a *authorizationService) GenerateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": fmt.Sprint(userID),
		"exp": time.Now().Add(time.Hour * 10).Unix(),
	})
	return token.SignedString([]byte(a.appConfig.JWT.Secret))
}

// ValidateJWT implements AuthorizationService.
func (a *authorizationService) ValidateJWT(token string) error {
	_, err := a.ParseJWTClaims(token)
	return err
}

// ValidateTelegramData implements AuthorizationService.
func (a *authorizationService) ValidateTelegramData(data string) error {
	rawToken := a.appConfig.TelegramBot.TokenNormalized
	return initdata.Validate(data, rawToken, time.Hour*24)
}

func New(appConfig *config.Config) AuthorizationService {
	return &authorizationService{
		appConfig: appConfig,
	}
}
