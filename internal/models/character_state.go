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
	Energy        uint8     `redis:"energy"`
	SpendEnergyAt time.Time `redis:"spend_energy_at"`
	State         State     `redis:"state"`
}
