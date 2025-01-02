package models

import "time"

type User struct {
	ID         int
	TelegramID int64
	UpdatedAt  time.Time
	CreatedAt  time.Time
}
