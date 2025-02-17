package resources

import "io"

type GameResoucesService interface {
	GetItems() io.Reader
	GetCrafts() io.Reader
	GetWeaponAttributes() io.Reader
	GetAdventures() io.Reader
	getAdventureRewards() io.Reader
}
