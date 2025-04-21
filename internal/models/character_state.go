package models

import (
	"errors"
	"time"
)

type State string

func (s State) MarshalBinary() (data []byte, err error) {
	return []byte(s), nil
}

func (s *State) UnmarshalBinary(data []byte) error {
	switch State(data) {
	case StateIdle, StateInAdventure, StateCrafting:
		*s = State(data)
		return nil
	default:
		return errors.New("invalid state value")
	}
}

const (
	StateIdle        State = "idle"
	StateInAdventure State = "adventure"
	StateCrafting    State = "crafting"
)

// CharacterState Any fast-updating non-important data store here
type CharacterState struct {
	Energy        uint8     `redis:"energy"`
	SpendEnergyAt time.Time `redis:"spend_energy_at"`
	State         State     `redis:"state"`
}
