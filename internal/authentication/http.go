package authentication

import (
	"log/slog"

	"github.com/yazmeyaa/duck_wars/internal/user"
)

type HTTPAuthService struct {
	log *slog.Logger
}

func (h *HTTPAuthService) Authenticate(token string) (user.User, error) {
	user := user.User{}

	return user, nil
}
