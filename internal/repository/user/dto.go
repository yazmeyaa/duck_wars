package user_repository

type CreateUserDTO struct {
	TelegramID int
	Username   string
}

type UpdateUserDTO struct {
	Username string
}
