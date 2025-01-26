package authorization_service

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"github.com/yazmeyaa/duck_wars/internal/config"
)

type AuthorizationService interface {
	ValidateTelegramData(data string) error
	ValidateJWT(token string) error
	GenerateToken(userID int) (string, error)
}

type authorizationService struct {
	appConfig *config.Config
}

// GenerateToken implements AuthorizationService.
func (a *authorizationService) GenerateToken(userID int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userID
	return token.SignedString([]byte(a.appConfig.JWT.Secret))
}

// ValidateJWT implements AuthorizationService.
func (a *authorizationService) ValidateJWT(token string) error {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(a.appConfig.JWT.Secret), nil
	})

	return err
}

// ValidateTelegramData implements AuthorizationService.
func (a *authorizationService) ValidateTelegramData(data string) error {
	rawToken := a.appConfig.TelegramBot.Token
	token := strings.TrimSuffix(rawToken, "/test")
	return initdata.Validate(data, token, time.Hour*24)
}

func New(appConfig *config.Config) AuthorizationService {
	return &authorizationService{
		appConfig: appConfig,
	}
}
