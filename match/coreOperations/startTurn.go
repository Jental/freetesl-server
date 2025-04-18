package coreOperations

import (
	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/match/models"
)

func StartTurn(playerState *models.PlayerMatchState) {
	if playerState.PlayerID == playerState.MatchState.PlayerWithFirstTurnID {
		playerState.MatchState.TurnID = playerState.MatchState.TurnID + 1
	}

	playerState.SetMaxMana(playerState.GetMaxMana() + 1)
	playerState.SetMana(playerState.GetMaxMana())

	if len(playerState.GetHand()) >= common.MAX_HAND_CARDS {
		DiscardCardFromDeck(playerState)
	} else {
		DrawCard(playerState)
	}

	playerState.SetRingActivity(playerState.GetRingGemCount() > 0)
}
