package user

import "time"

type User struct {
	ID         uint32
	Username   string
	Password   string
	TelegramID *int64
	Gmail      *string
	CreatedAt  time.Time
}
