package user

type CreateUserDTO struct {
	Username string
	Password string
}

type UpdateUserDTO struct {
	Username   *string
	Password   *string
	TelegramID *int64
	Gmail      *string
}
