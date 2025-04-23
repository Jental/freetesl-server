package models

import (
	"fmt"

	"github.com/google/uuid"
	dbEnums "github.com/jental/freetesl-server/db/enums"
	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/samber/lo"
)

type CardInstanceBase struct {
	Card             *dbModels.Card
	CardInstanceID   uuid.UUID
	Cost             int
	KeywordInstances []*KeywordInstance // all card types can have keywords, at least some, e.g. Prophecy
	IsActive         bool
}

type CardInstance interface {
	GetBase() *CardInstanceBase
	HasKeyword(keyword dbEnums.CardKeyword) bool
	IsActive() bool
	SetIsActive(isActive bool)
}

func newCardInstanceBase(card *dbModels.Card) CardInstanceBase {
	cardInstanceID := uuid.New()
	return CardInstanceBase{
		Card:           card,
		CardInstanceID: cardInstanceID,
		Cost:           card.Cost,
		KeywordInstances: lo.Map(card.Keywords, func(kw dbEnums.CardKeyword, _ int) *KeywordInstance {
			kwInst := NewKeywordInstance(kw, nil, cardInstanceID)
			return &kwInst
		}),
		IsActive: false,
	}
}

func NewCardInstance(card *dbModels.Card) (CardInstance, error) {
	switch card.Type {
	case dbEnums.CardTypeCreature:
		inst := NewCardInstanceCreature(card)
		return &inst, nil
	case dbEnums.CardTypeItem:
		inst, err := NewCardInstanceItem(card)
		if err != nil {
			return nil, err
		}
		return &inst, nil
	case dbEnums.CardTypeAction:
		inst := NewCardInstanceAction(card)
		return &inst, nil
	case dbEnums.CardTypeSupport:
		inst := NewCardInstanceSupport(card)
		return &inst, nil
	default:
		return nil, fmt.Errorf("NewCardInstance: Invalid card type: '%d'", card.Type)
	}
}

func cardInstanceHasKeyword(cardInstance CardInstance, keyword dbEnums.CardKeyword) bool {
	for _, kw := range cardInstance.GetBase().KeywordInstances {
		if kw.Keyword == keyword {
			return true
		}
	}
	return false
}

func cardInstanceIsActive(cardInstance CardInstance) bool {
	return cardInstance.GetBase().IsActive
}

func cardInstanceSetIsActive(cardInstance CardInstance, isActive bool) {
	cardInstance.GetBase().IsActive = isActive
}
