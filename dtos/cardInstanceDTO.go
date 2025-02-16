package dtos

import "github.com/google/uuid"

type CardInstanceDTO struct {
	CardInstanceID uuid.UUID `json:"cardInstanceID"`
	CardID         int       `json:"cardID"`
	Power          int       `json:"power"`
	Health         int       `json:"health"`
	Cost           int       `json:"cost"`
}
