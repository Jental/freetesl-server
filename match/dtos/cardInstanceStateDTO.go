package dtos

import "github.com/google/uuid"

type CardInstanceStateDTO struct {
	CardInstanceID uuid.UUID `json:"cardInstanceID"`
	IsActive       bool      `json:"isActive"`
}
