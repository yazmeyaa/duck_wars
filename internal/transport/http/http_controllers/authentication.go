package http_controllers

import (
	"net/http"
)

type AuthenticationController interface {
	LoginViaTelegramInitData(w http.ResponseWriter, r *http.Request)
	ValidateToken(w http.ResponseWriter, r *http.Request)
}

type authenticationController struct{}

// LoginViaTelegramInitData implements AuthenticationController.
func (a *authenticationController) LoginViaTelegramInitData(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// ValidateToken implements AuthenticationController.
func (a *authenticationController) ValidateToken(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func NewAuthenticationController() AuthenticationController {
	return &authenticationController{}
}
