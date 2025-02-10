package models

import "github.com/google/uuid"

type CardInstance struct {
	Card           Card
	CardInstanceID uuid.UUID
	Power          int
	Health         int
	Cost           int
	IsActive       bool
}
