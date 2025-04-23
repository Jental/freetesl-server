package models

import (
	"slices"

	dbEnums "github.com/jental/freetesl-server/db/enums"
	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/match/effects"
	"github.com/jental/freetesl-server/models/enums"
)

type CardInstanceCreature struct {
	CardInstanceBase
	power   int
	health  int
	Effects []*effects.EffectInstance
	Items   []*CardInstanceItem
}

func NewCardInstanceCreature(card *dbModels.Card) CardInstanceCreature {
	return CardInstanceCreature{
		CardInstanceBase: newCardInstanceBase(card),
		power:            card.Power,
		health:           card.Health,
		Effects:          make([]*effects.EffectInstance, 0),
		Items:            make([]*CardInstanceItem, 0),
	}
}

func (cardInstance *CardInstanceCreature) GetBase() *CardInstanceBase {
	return &cardInstance.CardInstanceBase
}

func (cardInstance *CardInstanceCreature) HasKeyword(keyword dbEnums.CardKeyword) bool {
	if cardInstanceHasKeyword(cardInstance, keyword) {
		return true
	}

	for _, item := range cardInstance.Items {
		if cardInstanceHasKeyword(item, keyword) {
			return true
		}
	}

	return false
}

func (cardInstance *CardInstanceCreature) GetAllKeywords() []*KeywordInstance {
	keywords := make([]*KeywordInstance, len(cardInstance.KeywordInstances))
	copy(keywords, cardInstance.KeywordInstances)
	for _, item := range cardInstance.Items {
		keywords = append(keywords, item.KeywordInstances...)
	}
	return keywords
}

func (cardInstance *CardInstanceCreature) IsActive() bool {
	return cardInstanceIsActive(cardInstance)
}

func (cardInstance *CardInstanceCreature) SetIsActive(isActive bool) {
	cardInstanceSetIsActive(cardInstance, isActive)
}

func (cardInstance *CardInstanceCreature) GetAllEffects() []*effects.EffectInstance {
	effs := make([]*effects.EffectInstance, len(cardInstance.Effects))
	copy(effs, cardInstance.Effects)
	for _, item := range cardInstance.Items {
		for _, eff := range item.Effects {
			startTurnID := -1
			if item.EquippedTurnID != nil {
				startTurnID = *item.EquippedTurnID
			}
			effectInstance := effects.NewEffectInstance(eff, startTurnID, &item.CardInstanceID)
			effs = append(effs, &effectInstance)
		}
	}
	return effs
}

func (cardInstance *CardInstanceCreature) HasEffect(effectType enums.EffectType) bool {
	idx := slices.IndexFunc(cardInstance.GetAllEffects(), func(eff *effects.EffectInstance) bool { return eff.Effect.GetType() == effectType })
	return idx >= 0
}

func (cardInstance *CardInstanceCreature) GetHealthIncrease() int {
	totalHealth := 0
	for _, item := range cardInstance.Items {
		for _, effect := range item.Effects {
			switch castedEffect := effect.(type) {
			case *effects.EffectModifyPowerHealth:
				totalHealth = totalHealth + castedEffect.HealthIncrease
			}
		}
	}
	return totalHealth
}

func (cardInstance *CardInstanceCreature) GetComputedHealth() int {
	return cardInstance.health + cardInstance.GetHealthIncrease()
}

func (cardIntance *CardInstanceCreature) UpdateHealth(updatedComputedHealth int) {
	cardIntance.health = updatedComputedHealth - cardIntance.GetHealthIncrease()
}

func (cardInstance *CardInstanceCreature) GetPowerIncrease() int {
	totalPower := 0
	for _, item := range cardInstance.Items {
		for _, effect := range item.Effects {
			switch castedEffect := effect.(type) {
			case *effects.EffectModifyPowerHealth:
				totalPower = totalPower + castedEffect.PowerIncrease
			}
		}
	}
	return totalPower
}

func (cardInstance *CardInstanceCreature) GetComputedPower() int {
	return cardInstance.power + cardInstance.GetPowerIncrease()
}

func (cardIntance *CardInstanceCreature) UpdatePower(updatedComputedPower int) {
	cardIntance.power = updatedComputedPower - cardIntance.GetPowerIncrease()
}
