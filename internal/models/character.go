package models

import "time"

type Character struct {
	ID         int
	UserID     int
	Name       string
	Coins      float64
	MaxEnergy  uint8
	Experience float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
