package dtos

type EffectInstanceDTO struct {
	ID                   byte    `json:"id"`
	Description          string  `json:"description"`
	SourceCardInstanceID *string `json:"sourceCardInstanceID"`
}
