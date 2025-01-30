package models

import "time"

type State string

const (
	STATE_IDLE         State = "idle"
	STATE_IN_ADVENTURE State = "adventure"
	STATE_CRAFTING     State = "crafting"
)

// Any fast-updating non-important data store here
type CharacterState struct {
	Energy        uint8
	SpendEnergyAt time.Time
	State         State
}
