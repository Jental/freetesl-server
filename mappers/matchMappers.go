package mappers

import (
	"fmt"

	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/models"
	"github.com/samber/lo"
)

func MapToCardInstanceDTO(model *models.CardInstance) dtos.CardInstanceDTO {
	return dtos.CardInstanceDTO{
		CardID:         model.Card.ID,
		CardInstanceID: model.CardInstanceID,
		IsActive:       model.IsActive,
		Power:          model.Power,
		Health:         model.Health,
		Cost:           model.Cost,
	}
}

func MapToPlayerMatchStateDTO(model *models.PlayerMatchState2) dtos.PlayerMatchStateDTO {
	return dtos.PlayerMatchStateDTO{
		Deck:           lo.Map(model.Deck, func(item *models.CardInstance, i int) dtos.CardInstanceDTO { return MapToCardInstanceDTO(item) }),
		Hand:           lo.Map(model.Hand, func(item *models.CardInstance, i int) dtos.CardInstanceDTO { return MapToCardInstanceDTO(item) }),
		Health:         model.Health,
		Runes:          model.Runes,
		Mana:           model.Mana,
		MaxMana:        model.MaxMana,
		LeftLaneCards:  lo.Map(model.LeftLaneCards, func(item *models.CardInstance, i int) dtos.CardInstanceDTO { return MapToCardInstanceDTO(item) }),
		RightLaneCards: lo.Map(model.RightLaneCards, func(item *models.CardInstance, i int) dtos.CardInstanceDTO { return MapToCardInstanceDTO(item) }),
	}
}

func MapToMatchStateDTO(model *models.Match, playerID int) (*dtos.MatchStateDTO, error) {
	var playerState *models.PlayerMatchState2
	var opponentState *models.PlayerMatchState2 = nil
	if model.Player0State.HasValue && model.Player0State.Value.PlayerID == playerID {
		playerState = model.Player0State.Value
		if model.Player1State.HasValue {
			opponentState = model.Player1State.Value
		}
	} else if model.Player1State.HasValue && model.Player1State.Value.PlayerID == playerID {
		playerState = model.Player1State.Value
		if model.Player0State.HasValue {
			opponentState = model.Player0State.Value
		}
	} else {
		return nil, fmt.Errorf("player with id '%d' is not a part of a match", playerID)
	}

	return &dtos.MatchStateDTO{
		Player:   MapToPlayerMatchStateDTO(playerState),
		Opponent: MapToPlayerMatchStateDTO(opponentState),
		OwnTurn:  model.PlayerWithTurnID == playerID,
	}, nil
}

func MapToMatchEndDTO(model *models.Match, playerID int) dtos.MatchEndDTO {
	return dtos.MatchEndDTO{
		HasWon: model.WinnerID == playerID,
	}
}
