package http_controllers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	user_repository "github.com/yazmeyaa/ton_game/internal/repository/user"
	authorization_service "github.com/yazmeyaa/ton_game/internal/service/authorization"
	user_service "github.com/yazmeyaa/ton_game/internal/service/user"
	"github.com/yazmeyaa/ton_game/internal/transport/transport_http"
)

type AuthenticationController interface {
	LoginViaTelegramInitData(w http.ResponseWriter, r *http.Request)
	ValidateToken(w http.ResponseWriter, r *http.Request)
}

type AuthResponse struct {
	Token string `json:"token"`
}

type authenticationController struct {
	authSvc authorization_service.AuthorizationService
	userSvc user_service.UserService
	logger  *slog.Logger
}

func (a *authenticationController) writeToken(w http.ResponseWriter, userID int) {
	token, err := a.authSvc.GenerateToken(userID)
	if err != nil {
		err := transport_http.WriteError(w, http.StatusInternalServerError, "Internal server error. Cannot generate token.")
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := AuthResponse{
		Token: token,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return
	}
}

// LoginViaTelegramInitData implements AuthenticationController.
func (a *authenticationController) LoginViaTelegramInitData(w http.ResponseWriter, r *http.Request) {
	const fn string = "AuthenticationController.LoginViaTelegramInitData"
	type requestBody struct {
		InitData string `json:"initData"`
	}
	var requestData requestBody
	if err := transport_http.ParseJSON(&requestData, w, r, a.logger, fn); err != nil {
		return
	}

	if err := a.authSvc.ValidateTelegramData(requestData.InitData); err != nil {
		a.logger.Error("Error while validating telegram init data", "fn", fn, "error", err.Error())
		err := transport_http.WriteError(w, http.StatusUnauthorized, "Not valid token")
		if err != nil {
			return
		}
		return
	}

	data, err := initdata.Parse(requestData.InitData)
	if err != nil {
		a.logger.Error("Cannot parse telegram init data", "fn", fn, "error", err.Error())
		err := transport_http.WriteError(w, http.StatusInternalServerError, "Internal server error")
		if err != nil {
			return
		}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	user, err := a.userSvc.GetByTelegramID(ctx, data.User.ID)

	var userID int
	if errors.Is(err, user_repository.ErrUserNotFound) {
		newUserID, err := a.userSvc.Create(ctx, user_service.CreateUserDTO{
			TelegramID: data.User.ID,
		})
		if err != nil {
			err := transport_http.WriteError(w, http.StatusInternalServerError, "Cannot create user.")
			if err != nil {
				return
			}
			return
		}
		userID = newUserID
		a.writeToken(w, newUserID)
		return
	} else if err != nil {
		err := transport_http.WriteError(w, http.StatusInternalServerError, "Cannot get or create user...")
		if err != nil {
			return
		}
		return
	}
	userID = user.ID

	a.writeToken(w, userID)
}

func (a *authenticationController) processValidationError(w http.ResponseWriter, err error, logger *slog.Logger) {
	const fn = "processValidationError"
	var status int
	var message string

	if err != nil {
		logger.Error("Error while validating JWT token.", "fn", fn, "error", err.Error())
	}

	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		status = http.StatusUnauthorized
		message = "Token expired"
	case errors.Is(err, jwt.ErrTokenMalformed):
		status = http.StatusUnauthorized
		message = "Malformed token"
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		status = http.StatusUnauthorized
		message = "Invalid token signature"
	default:
		status = http.StatusInternalServerError
		message = "Internal server error. Cannot verify token"
	}

	err = transport_http.WriteError(w, status, message)
	if err != nil {
		return
	}
}

// ValidateToken implements AuthenticationController.
func (a *authenticationController) ValidateToken(w http.ResponseWriter, r *http.Request) {
	const fn = "AuthenticationController.ValidateToken"
	type requestBody struct {
		Token string `json:"token"`
	}
	var body requestBody

	if err := transport_http.ParseJSON(&body, w, r, a.logger, fn); err != nil {
		return
	}

	if body.Token == "" {
		a.logger.Error("Received empty token.", "fn", fn)
		err := transport_http.WriteError(w, http.StatusBadRequest, "Received empty token.")
		if err != nil {
			return
		}
		return
	}

	err := a.authSvc.ValidateJWT(body.Token)
	if err != nil {
		a.processValidationError(w, err, a.logger)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewAuthenticationController(authSvc authorization_service.AuthorizationService, userSvc user_service.UserService, logger *slog.Logger) AuthenticationController {
	return &authenticationController{
		authSvc: authSvc,
		userSvc: userSvc,
		logger:  logger,
	}
}
