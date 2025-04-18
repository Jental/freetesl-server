package models

import (
	dbEnums "github.com/jental/freetesl-server/db/enums"
	dbModels "github.com/jental/freetesl-server/db/models"
)

type CardInstanceSupport struct {
	CardInstanceBase
}

func NewCardInstanceSupport(card *dbModels.Card) CardInstanceSupport {
	return CardInstanceSupport{
		CardInstanceBase: newCardInstanceBase(card),
	}
}

func (cardInstance *CardInstanceSupport) GetBase() *CardInstanceBase {
	return &cardInstance.CardInstanceBase
}

func (cardInstance *CardInstanceSupport) HasKeyword(keyword dbEnums.CardKeyword) bool {
	return cardInstanceHasKeyword(cardInstance, keyword)
}

func (cardInstance *CardInstanceSupport) IsActive() bool {
	return cardInstanceIsActive(cardInstance)
}

func (cardInstance *CardInstanceSupport) SetIsActive(isActive bool) {
	cardInstanceSetIsActive(cardInstance, isActive)
}
