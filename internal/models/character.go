package models

type Character struct {
	ID        int
	UserID    int
	Name      string
	Coins     float64
	MaxEnergy uint8
}
