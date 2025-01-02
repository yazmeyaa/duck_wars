package game_data

type WeaponAttribute struct {
	ItemID         string          `json:"item_id"`
	Attack         AttackAttribute `json:"attack"`
	CriticalChance float64         `json:"critical_chacne"`
	AttackSpeed    float64         `json:"attack_speed"`
}

type AttackAttribute struct {
	Min int `json:"min"`
	Max int `json:"max"`
}
