package coreOperations

import "github.com/jental/freetesl-server/models"

func StartTurn(playerState *models.PlayerMatchState) {
	if playerState.PlayerID == playerState.MatchState.PlayerWithFirstTurnID {
		playerState.MatchState.TurnID = playerState.MatchState.TurnID + 1
	}

	playerState.SetMaxMana(playerState.GetMaxMana() + 1)
	playerState.SetMana(playerState.GetMaxMana())
	DrawCard(playerState)
}
