package mappers

import (
	"fmt"

	dbEnums "github.com/jental/freetesl-server/db/enums"
	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/models"
	"github.com/samber/lo"
)

func MapToCardInstanceDTO(model *models.CardInstance) dtos.CardInstanceDTO {
	return dtos.CardInstanceDTO{
		CardID:         model.Card.ID,
		CardInstanceID: model.CardInstanceID,
		Power:          model.Power,
		Health:         model.Health,
		Cost:           model.Cost,
		Keywords:       lo.Map(model.Keywords, func(kwd dbEnums.CardKeyword, _ int) int { return int(kwd) }),
		Effects:        lo.Map(model.Effects, func(eff *models.Effect, _ int) int { return int(eff.EffectType) }),
		// TODO:
		// - send unique effect types
		// - some effects (like silence) may overlap other effects - send only ones actual for FE
	}
}

func MapToCardInstanceStateDTO(model *models.CardInstance) dtos.CardInstanceStateDTO {
	return dtos.CardInstanceStateDTO{
		CardInstanceID: model.CardInstanceID,
		IsActive:       model.IsActive,
	}
}

func MapToPlayerMatchStateDTO(model *models.PlayerMatchState) dtos.PlayerMatchStateDTO {
	result := dtos.PlayerMatchStateDTO{
		Health:  model.GetHealth(),
		Runes:   model.GetRunes(),
		Mana:    model.GetMana(),
		MaxMana: model.GetMaxMana(),
		Hand: lo.Map(model.GetHand(), func(item *models.CardInstance, i int) dtos.CardInstanceStateDTO {
			return MapToCardInstanceStateDTO(item)
		}),
		LeftLaneCards: lo.Map(model.GetLeftLaneCards(), func(item *models.CardInstance, i int) dtos.CardInstanceStateDTO {
			return MapToCardInstanceStateDTO(item)
		}),
		RightLaneCards: lo.Map(model.GetRightLaneCards(), func(item *models.CardInstance, i int) dtos.CardInstanceStateDTO {
			return MapToCardInstanceStateDTO(item)
		}),
	}

	cardInstanceForAction := model.GetCardInstanceWaitingForAction()
	if cardInstanceForAction != nil {
		result.CardInstanceWaitingForAction = &cardInstanceForAction.CardInstanceID
	}

	return result
}

func MapToMatchStateDTO(model *models.Match, playerID int) (*dtos.MatchStateDTO, error) {
	var playerState *models.PlayerMatchState
	var opponentState *models.PlayerMatchState = nil
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

func MapToDeckStateDTO(model *models.Match, playerID int) (*dtos.DeckStateDTO, error) {
	var playerState *models.PlayerMatchState
	var opponentState *models.PlayerMatchState = nil
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

	return &dtos.DeckStateDTO{
		Player: lo.Map(playerState.GetDeck(), func(item *models.CardInstance, i int) *dtos.CardInstanceStateDTO {
			var r = MapToCardInstanceStateDTO(item)
			return &r
		}),
		Opponent: lo.Map(opponentState.GetDeck(), func(item *models.CardInstance, i int) *dtos.CardInstanceStateDTO {
			var r = MapToCardInstanceStateDTO(item)
			return &r
		}),
	}, nil
}

func MapToDiscardPileStateDTO(model *models.Match, playerID int) (*dtos.DiscardPileStateDTO, error) {
	var playerState *models.PlayerMatchState
	var opponentState *models.PlayerMatchState = nil
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

	return &dtos.DiscardPileStateDTO{
		Player: lo.Map(playerState.GetDiscardPile(), func(item *models.CardInstance, i int) *dtos.CardInstanceStateDTO {
			var r = MapToCardInstanceStateDTO(item)
			return &r
		}),
		Opponent: lo.Map(opponentState.GetDiscardPile(), func(item *models.CardInstance, i int) *dtos.CardInstanceStateDTO {
			var r = MapToCardInstanceStateDTO(item)
			return &r
		}),
	}, nil
}

func MapToMatchEndDTO(model *models.Match, playerID int) dtos.MatchEndDTO {
	return dtos.MatchEndDTO{
		HasWon: model.WinnerID == playerID,
	}
}
