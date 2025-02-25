package handlers

import (
	"fmt"
	"time"

	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/actions"
	"github.com/jental/freetesl-server/models"
)

func EndTurn(playerID int) {
	matchState, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Println(err)
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

// TODO: remove after bot implemented
func EndTurnAuto(playerID int) {
	matchState, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Println(err)
		return
	}

	actions.SwitchTurn(matchState)
	opponentState.SetMaxMana(opponentState.GetMaxMana() + 1)
	opponentState.SetMana(opponentState.GetMaxMana())
	time.Sleep(1 * time.Second)
	actions.DrawCard(opponentState)
	actions.PlayRandomCards(opponentState, playerState, matchState)

	for _, card := range opponentState.GetLeftLaneCards() {
		card.IsActive = false
	}
	for _, card := range opponentState.GetRightLaneCards() {
		card.IsActive = false
	}

	actions.SwitchTurn(matchState)
	actions.DrawCard(playerState)

	for _, card := range playerState.GetLeftLaneCards() {
		card.IsActive = true
	}
	for _, card := range playerState.GetRightLaneCards() {
		card.IsActive = true
	}

	playerState.SetMaxMana(playerState.GetMaxMana() + 1)
	playerState.SetMana(playerState.GetMaxMana())
}
