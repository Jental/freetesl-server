package mappers

import (
	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/models"
	"github.com/samber/lo"
)

func MapToCardInstanceDTO(model *models.CardInstance) dtos.CardInstanceDTO {
	return dtos.CardInstanceDTO{
		CardID:         model.Card.ID,
		CardInstanceID: model.CardInstanceID,
	}
}

func MapToPlayerMatchStateDTO(model *models.PlayerMatchState2, ownTurn bool) dtos.PlayerMatchStateDTO {
	return dtos.PlayerMatchStateDTO{
		Deck:    lo.Map(model.Deck, func(item *models.CardInstance, i int) dtos.CardInstanceDTO { return MapToCardInstanceDTO(item) }),
		Hand:    lo.Map(model.Hand, func(item *models.CardInstance, i int) dtos.CardInstanceDTO { return MapToCardInstanceDTO(item) }),
		Health:  model.Health,
		Runes:   model.Runes,
		Mana:    model.Mana,
		MaxMana: model.MaxMana,
		OwnTurn: ownTurn,
	}
}
