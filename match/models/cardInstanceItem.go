package models

import (
	"slices"

	dbEnums "github.com/jental/freetesl-server/db/enums"
	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/models/enums"
)

type CardInstanceItem struct {
	CardInstanceBase
	PowerIncrease  int
	HealthIncrease int
	Effects        []*Effect
}

func NewCardInstanceItem(card *dbModels.Card) CardInstanceItem {
	return CardInstanceItem{
		CardInstanceBase: newCardInstanceBase(card),
		PowerIncrease:    card.Power,
		HealthIncrease:   card.Health,
		Effects:          make([]*Effect, 0),
	}
}

func (cardInstance *CardInstanceItem) GetBase() *CardInstanceBase {
	return &cardInstance.CardInstanceBase
}

func (cardInstance *CardInstanceItem) HasKeyword(keyword dbEnums.CardKeyword) bool {
	return cardInstanceHasKeyword(cardInstance, keyword)
}

func (cardInstance *CardInstanceItem) IsActive() bool {
	return cardInstanceIsActive(cardInstance)
}

func (cardInstance *CardInstanceItem) SetIsActive(isActive bool) {
	cardInstanceSetIsActive(cardInstance, isActive)
}

func (cardInstance *CardInstanceItem) HasEffect(effectType enums.EffectType) bool {
	idx := slices.IndexFunc(cardInstance.Effects, func(eff *Effect) bool { return eff.EffectType == effectType })
	return idx >= 0
}
