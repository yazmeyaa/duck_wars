package user

import "log/slog"

type UserService interface {
	Create(data *CreateUserDTO) (*User, error)
	Delete(user *User) error
	DeleteById(id uint32) error
	Update(user *User, data *UpdateUserDTO) error
	GetByUsername(username string) (*User, error)
	GetById(id uint32) (*User, error)
}

type userService struct {
	userRepository UserRepository
	log            *slog.Logger
}

func (u *userService) Create(data *CreateUserDTO) (*User, error) {
	return u.userRepository.Create(data)
}

func (u *userService) Delete(user *User) error {
	return u.userRepository.Delete(user)
}

func (u *userService) DeleteById(id uint32) error {
	return u.userRepository.DeleteById(id)
}

func (u *userService) GetById(id uint32) (*User, error) {
	return u.userRepository.GetById(id)
}

func (u *userService) GetByUsername(username string) (*User, error) {
	return u.userRepository.GetByUsername(username)
}

func (u *userService) Update(user *User, data *UpdateUserDTO) error {
	return u.userRepository.Update(user, data)
}

func NewUserSerivce(userRepository UserRepository, log *slog.Logger) UserService {
	return &userService{
		userRepository: userRepository,
		log:            log,
	}
}
