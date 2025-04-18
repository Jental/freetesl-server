package handlers

import (
	"fmt"

	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/match/match"
	"github.com/jental/freetesl-server/match/operations"
)

func EndTurn(playerID int) {
	matchState, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		return
	}

	for _, card := range playerState.GetLeftLaneCards() {
		card.SetIsActive(false)
	}
	for _, card := range playerState.GetRightLaneCards() {
		card.SetIsActive(false)
	}

	coreOperations.SwitchTurn(matchState)
	operations.StartTurn(opponentState, matchState)
}
