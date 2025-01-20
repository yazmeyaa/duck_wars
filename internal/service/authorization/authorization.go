package authorization_service

import (
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
	"github.com/yazmeyaa/duck_wars/internal/config"
)

type AuthorizationService interface {
	ValidateTelegramData(data string) error
	ValidateJWT(token string) (bool, error)
}

type authorizationService struct {
	appConfig *config.Config
}

// ValidateJWT implements AuthorizationService.
func (a *authorizationService) ValidateJWT(token string) (bool, error) {
	panic("unimplemented")
}

// ValidateTelegramData implements AuthorizationService.
func (a *authorizationService) ValidateTelegramData(data string) error {
	return initdata.Validate(data, a.appConfig.TelegramBot.Token, time.Hour*24)
}

func New(appConfig *config.Config) AuthorizationService {
	return &authorizationService{
		appConfig: appConfig,
	}
}
