package coreOperations

import (
	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/models/enums"
)

func DrawCard(playerState *models.PlayerMatchState) models.CardInstance {
	var deck = playerState.GetDeck()
	if len(deck) == 0 {
		return nil // TODO: for now doing nothing, but later next rune should be broken
	}

	var drawnCard = deck[0]
	drawnCard.SetIsActive(true)
	playerState.SetHand(append(playerState.GetHand(), drawnCard))
	playerState.SetDeck(deck[1:])

	return drawnCard
}

func SwitchTurn(matchState *models.Match) {
	var isFirstPlayersTurn = matchState.Player0State.Value.PlayerID == matchState.PlayerWithTurnID
	if isFirstPlayersTurn {
		matchState.PlayerWithTurnID = matchState.Player1State.Value.PlayerID
	} else {
		matchState.PlayerWithTurnID = matchState.Player0State.Value.PlayerID
	}

	matchState.Player0State.Value.SendEvent(enums.BackendEventSwitchTurn)
	matchState.Player1State.Value.SendEvent(enums.BackendEventSwitchTurn)
}

func ReducePlayerHealth(playerState *models.PlayerMatchState, amount int) {
	var updatedHealth = playerState.GetHealth() - amount
	playerState.SetHealth(updatedHealth)

	var expectedRuneCount uint8 = uint8((updatedHealth - 1) / 5)
	var runeCount = max(0, min(expectedRuneCount, playerState.GetRunes()))
	playerState.SetRunes(runeCount)
}

func IncreasePlayerHealth(playerState *models.PlayerMatchState, amount int) {
	var updatedHealth = playerState.GetHealth() + amount
	playerState.SetHealth(updatedHealth)
}

func ReduceCardHealth(playerState *models.PlayerMatchState, cardInstance *models.CardInstanceCreature, lane *models.Lane, amount int) {
	updatedHealth := cardInstance.GetComputedHealth() - amount
	cardInstance.UpdateHealth(updatedHealth)
	if updatedHealth <= 0 {
		DiscardCardFromLane(playerState, cardInstance, lane)
	}

	playerState.SendEvent(enums.BackendEventCardInstancesChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentCardInstancesChanged)

	// to force lanes redraw
	playerState.SendEvent(enums.BackendEventLanesChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentLanesChanged)
}
