package user_service

import (
	"context"
	"log/slog"

	"github.com/yazmeyaa/duck_wars/internal/models"
	user_repository "github.com/yazmeyaa/duck_wars/internal/repository/user"
)

type UserService interface {
	Create(ctx context.Context, data user_repository.CreateUserDTO) (int, error)
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByTelegramID(ctx context.Context, id int64) (*models.User, error)
}

type userService struct {
	logger         *slog.Logger
	userRepository user_repository.UserRepository
}

// Create implements UserService.
func (u *userService) Create(ctx context.Context, data user_repository.CreateUserDTO) (int, error) {
	id, err := u.userRepository.Create(ctx, data)
	if err != nil {
		u.logger.Error("Cannot create user.", "error", err.Error())
		return 0, err
	}
	return id, nil
}

// Delete implements UserService.
func (u *userService) Delete(ctx context.Context, id int) error {
	err := u.userRepository.Delete(ctx, id)
	if err != nil {
		u.logger.Error("Cannot delete user.", "error", err.Error())
		return err
	}
	return nil
}

// GetByID implements UserService.
func (u *userService) GetByID(ctx context.Context, id int) (*models.User, error) {
	user, err := u.userRepository.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Cannot get user.", "error", err.Error())
		return nil, err
	}

	return user, nil
}

// GetByTelegramID implements UserService.
func (u *userService) GetByTelegramID(ctx context.Context, id int64) (*models.User, error) {
	user, err := u.userRepository.GetByTelegramID(ctx, id)
	if err != nil {
		u.logger.Error("Cannot get user.", "error", err.Error())
		return nil, err
	}

	return user, nil
}

func New(logger *slog.Logger, userRepository user_repository.UserRepository) UserService {
	return &userService{logger: logger, userRepository: userRepository}
}
