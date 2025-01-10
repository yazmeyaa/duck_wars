package game_data

type AdventureTaskPayload struct {
	AdventureID string `json:"adventureId"`
	UserID      int64  `json:"userId"`
	CharacterID int    `json:"characterId"`
}
