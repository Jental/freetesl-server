package models

import (
	"github.com/google/uuid"
	"github.com/jental/freetesl-server/db/enums"
	"github.com/jental/freetesl-server/db/models"
)

type CardInstance struct {
	Card           *models.Card
	CardInstanceID uuid.UUID
	Power          int
	Health         int
	Cost           int
	Keywords       []enums.CardKeyword
	IsActive       bool
}
