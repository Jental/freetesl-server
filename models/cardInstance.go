package models

import (
	"slices"

	"github.com/google/uuid"
	dbEnums "github.com/jental/freetesl-server/db/enums"
	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/models/enums"
)

type CardInstance struct {
	Card           *dbModels.Card
	CardInstanceID uuid.UUID
	Power          int
	Health         int
	Cost           int
	Keywords       []dbEnums.CardKeyword
	IsActive       bool
	Effects        []*Effect
}

func (cardInstance *CardInstance) HasEffect(effectType enums.EffectType) bool {
	idx := slices.IndexFunc(cardInstance.Effects, func(eff *Effect) bool { return eff.EffectType == effectType })
	return idx >= 0
}
