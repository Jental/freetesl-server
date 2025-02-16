package dtos

type PlayerMatchStateDTO struct {
	Health         int                    `json:"health"`
	Runes          uint8                  `json:"runes"`
	Mana           int                    `json:"mana"`
	MaxMana        int                    `json:"maxMana"`
	Hand           []CardInstanceStateDTO `json:"hand"`
	LeftLaneCards  []CardInstanceStateDTO `json:"leftLaneCards"`
	RightLaneCards []CardInstanceStateDTO `json:"rightLaneCards"`
}
