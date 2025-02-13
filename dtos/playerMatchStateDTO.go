package dtos

type PlayerMatchStateDTO struct {
	Deck           []CardInstanceDTO `json:"deck"`
	Hand           []CardInstanceDTO `json:"hand"`
	Health         int               `json:"health"`
	Runes          uint8             `json:"runes"`
	Mana           int               `json:"mana"`
	MaxMana        int               `json:"maxMana"`
	LeftLaneCards  []CardInstanceDTO `json:"leftLaneCards"`
	RightLaneCards []CardInstanceDTO `json:"rightLaneCards"`
	DiscardPile    []CardInstanceDTO `json:"discardPile"`
}
