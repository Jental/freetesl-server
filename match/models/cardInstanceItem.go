package models

import (
	"github.com/jental/freetesl-server/common"
	dbEnums "github.com/jental/freetesl-server/db/enums"
	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/match/effects"
	"github.com/jental/freetesl-server/models/enums"
	"github.com/samber/lo"
)

type CardInstanceItem struct {
	CardInstanceBase
	Effects []effects.IEffect
}

func NewCardInstanceItem(card *dbModels.Card) (CardInstanceItem, error) {
	effects, err := common.MapErr(card.Effects, func(dbEffect dbModels.CardEffect, _ int) (effects.IEffect, error) {
		return effects.NewEffect(dbEffect)
	})
	if err != nil {
		return CardInstanceItem{}, err
	}

	return CardInstanceItem{
		CardInstanceBase: newCardInstanceBase(card),
		Effects:          effects,
	}, nil
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
	found := cardInstance.GetEffects(effectType)
	return len(found) > 0
}

func (cardInstance *CardInstanceItem) GetEffects(effectType enums.EffectType) []effects.IEffect {
	return lo.Filter(cardInstance.Effects, func(eff effects.IEffect, _ int) bool { return eff.GetType() == effectType })
}
