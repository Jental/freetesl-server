package dtos

type PlayerMatchStateDTO struct {
	Deck           []CardInstanceStateDTO `json:"deck"`
	Hand           []CardInstanceStateDTO `json:"hand"`
	Health         int                    `json:"health"`
	Runes          uint8                  `json:"runes"`
	Mana           int                    `json:"mana"`
	MaxMana        int                    `json:"maxMana"`
	LeftLaneCards  []CardInstanceStateDTO `json:"leftLaneCards"`
	RightLaneCards []CardInstanceStateDTO `json:"rightLaneCards"`
	DiscardPile    []CardInstanceStateDTO `json:"discardPile"`
}
