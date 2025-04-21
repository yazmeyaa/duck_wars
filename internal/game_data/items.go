package game_data

type Item struct {
	ID            string   `json:"id"`
	Type          ItemType `json:"type"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	CanSold       bool     `json:"canSold"`
	IsCraftable   bool     `json:"craftable"`
	IsUsedInCraft bool     `json:"usedInCraft"`
}
