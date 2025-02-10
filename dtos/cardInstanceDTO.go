package dtos

import "github.com/google/uuid"

type CardInstanceDTO struct {
	CardID         int       `json:"cardID"`
	CardInstanceID uuid.UUID `json:"cardInstanceID"`
	IsActive       bool      `json:"isActive"`
}
