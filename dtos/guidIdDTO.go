package dtos

import "github.com/google/uuid"

type GuidIdDTO struct {
	ID *uuid.UUID `json:"id"`
}
