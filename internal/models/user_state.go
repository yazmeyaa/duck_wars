package models

type State string

const (
	STATE_IDLE      State = "idle"
	STATE_TRAVELING State = "traveling"
)

// Any fast-updating not important data store here
type CharacterState struct {
	Energy uint8
	State  State
}
