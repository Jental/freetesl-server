package handlers

import (
	"fmt"

	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/coreOperations"
)

func EndTurn(playerID int) {
	matchState, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		return
	}

	for _, card := range playerState.GetLeftLaneCards() {
		card.IsActive = false
	}
	for _, card := range playerState.GetRightLaneCards() {
		card.IsActive = false
	}

	coreOperations.SwitchTurn(matchState)
	coreOperations.StartTurn(opponentState)
}
