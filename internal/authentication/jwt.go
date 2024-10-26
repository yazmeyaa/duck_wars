package authentication

import (
	"fmt"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yazmeyaa/duck_wars/config"
	"github.com/yazmeyaa/duck_wars/internal/user"
	"github.com/yazmeyaa/duck_wars/pkg/logger"
)

type JWTService struct {
	log    *slog.Logger
	config *config.Config
}

func (j *JWTService) CreateToken(user *user.User) (string, error) {
	const funcName string = "JWTService.CreateToken"

	claims := jwt.MapClaims{
		"sub": user.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(j.config.JWT.Secret)
	if err != nil {
		logger.LogError(j.log, funcName, "Failed to sign token", err)
		return "", err
	}
	return signed, err
}

func (j *JWTService) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return j.config.JWT.Secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
