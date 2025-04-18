package models

import (
	"slices"

	dbEnums "github.com/jental/freetesl-server/db/enums"
	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/models/enums"
)

type CardInstanceCreature struct {
	CardInstanceBase
	Power   int
	Health  int
	Effects []*Effect
	Items   []*CardInstanceItem
}

func NewCardInstanceCreature(card *dbModels.Card) CardInstanceCreature {
	return CardInstanceCreature{
		CardInstanceBase: newCardInstanceBase(card),
		Power:            card.Power,
		Health:           card.Health,
		Effects:          make([]*Effect, 0),
		Items:            make([]*CardInstanceItem, 0),
	}
}

func (cardInstance *CardInstanceCreature) GetBase() *CardInstanceBase {
	return &cardInstance.CardInstanceBase
}

func (cardInstance *CardInstanceCreature) HasKeyword(keyword dbEnums.CardKeyword) bool {
	return cardInstanceHasKeyword(cardInstance, keyword)
}

func (cardInstance *CardInstanceCreature) IsActive() bool {
	return cardInstanceIsActive(cardInstance)
}

func (cardInstance *CardInstanceCreature) SetIsActive(isActive bool) {
	cardInstanceSetIsActive(cardInstance, isActive)
}

func (cardInstance *CardInstanceCreature) HasEffect(effectType enums.EffectType) bool {
	idx := slices.IndexFunc(cardInstance.Effects, func(eff *Effect) bool { return eff.EffectType == effectType })
	return idx >= 0
}
