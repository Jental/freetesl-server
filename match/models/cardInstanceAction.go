package models

import (
	dbEnums "github.com/jental/freetesl-server/db/enums"
	dbModels "github.com/jental/freetesl-server/db/models"
)

type CardInstanceAction struct {
	CardInstanceBase
}

func NewCardInstanceAction(card *dbModels.Card) CardInstanceAction {
	return CardInstanceAction{
		CardInstanceBase: newCardInstanceBase(card),
	}
}

func (cardInstance *CardInstanceAction) GetBase() *CardInstanceBase {
	return &cardInstance.CardInstanceBase
}

func (cardInstance *CardInstanceAction) HasKeyword(keyword dbEnums.CardKeyword) bool {
	return cardInstanceHasKeyword(cardInstance, keyword)
}

func (cardInstance *CardInstanceAction) IsActive() bool {
	return cardInstanceIsActive(cardInstance)
}

func (cardInstance *CardInstanceAction) SetIsActive(isActive bool) {
	cardInstanceSetIsActive(cardInstance, isActive)
}
