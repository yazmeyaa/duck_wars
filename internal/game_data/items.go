package game_data

type Item struct {
	ID            string   `json:"id"`
	Type          ItemType `json:"type"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	IsSoldable    bool     `json:"soldable"`
	IsCraftable   bool     `json:"craftable"`
	IsUsedInCraft bool     `json:"userInCraft"`
}
