package game_data

type Adventure struct {
	AdventureID     string   `json:"adventure_id"`
	Title           string   `json:"title"`
	SuccessChance   float32  `json:"success_chance"`
	Description     string   `json:"description"`
	Duration        int      `json:"duration"`
	EnergyCost      int      `json:"energy_cost"`
	StartMessages   []string `json:"start_messages"`
	SuccessMessages []string `json:"success_messages"`
	FailureMessages []string `json:"failure_messages"`
}

type AdventureRewards struct {
	AdventureID string   `json:"adventure_id"`
	Experience  MinMax   `json:"experience"`
	Rewards     []Reward `json:"rewards"`
}

type MinMax struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type Reward struct {
	ItemID   string  `json:"item_id"`
	Chance   float32 `json:"chance"`
	Quantity MinMax  `json:"quantity"`
}
