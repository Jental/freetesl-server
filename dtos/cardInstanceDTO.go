package dtos

import "github.com/google/uuid"

type CardInstanceDTO struct {
	CardInstanceID uuid.UUID `json:"cardInstanceID"`
	CardID         int       `json:"cardID"`
	Power          int       `json:"power"`
	Health         int       `json:"health"`
	Cost           int       `json:"cost"`
	Keywords       []int     `json:"keywords"` // we don't use []byte here because it's serialized as Base64
}
