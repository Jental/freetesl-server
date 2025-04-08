package dtos

import "github.com/google/uuid"

type PlayerMatchStateDTO struct {
	Health                       int                    `json:"health"`
	Runes                        uint8                  `json:"runes"`
	Mana                         int                    `json:"mana"`
	MaxMana                      int                    `json:"maxMana"`
	Hand                         []CardInstanceStateDTO `json:"hand"`
	LeftLaneCards                []CardInstanceStateDTO `json:"leftLaneCards"`
	RightLaneCards               []CardInstanceStateDTO `json:"rightLaneCards"`
	RingGemCount                 uint8                  `json:"ringGemCount"`
	IsRingActive                 bool                   `json:"isRingActive"`
	CardInstanceWaitingForAction *uuid.UUID             `json:"cardInstanceWaitingForAction"`
}
