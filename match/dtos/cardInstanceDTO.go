package dtos

import "github.com/google/uuid"

type CardInstanceDTO struct {
	CardInstanceID uuid.UUID            `json:"cardInstanceID"`
	CardID         int                  `json:"cardID"`
	Power          int                  `json:"power"`     // modification is included
	PowerMod       int                  `json:"powerMod"`  // 0 - not modified; <0 - decreased; >0 - increased
	Health         int                  `json:"health"`    // modification is included
	HealthMod      int                  `json:"healthMod"` // 0 - not modified; <0 - decreased; >0 - increased
	Cost           int                  `json:"cost"`
	Keywords       []KeywordInstanceDTO `json:"keywords"`
	Effects        []EffectInstanceDTO  `json:"effects"`
}
