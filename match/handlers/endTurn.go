package handlers

import (
	"fmt"

	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/actions"
	"github.com/jental/freetesl-server/models"
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

	actions.SwitchTurn(matchState)

	startTurn(opponentState)
}

func startTurn(playerState *models.PlayerMatchState) {
	playerState.SetMaxMana(playerState.GetMaxMana() + 1)
	playerState.SetMana(playerState.GetMaxMana())
	actions.DrawCard(playerState)

	for _, card := range playerState.GetLeftLaneCards() {
		card.IsActive = true
	}
	for _, card := range playerState.GetRightLaneCards() {
		card.IsActive = true
	}
}
