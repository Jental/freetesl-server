package coreOperations

import (
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

func DrawCard(playerState *models.PlayerMatchState) {
	var deck = playerState.GetDeck()
	if len(deck) == 0 {
		return // TODO: for now doing nothing, but later next rune should be broken
	}

	var drawnCard = deck[0]
	drawnCard.IsActive = true
	playerState.SetHand(append(playerState.GetHand(), drawnCard))
	playerState.SetDeck(deck[1:])
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

func PlaceCardToLane(playerState *models.PlayerMatchState, cardInstance *models.CardInstance, laneID enums.Lane) {
	if laneID == enums.LaneLeft {
		playerState.SetLeftLaneCards(append(playerState.GetLeftLaneCards(), cardInstance))
	} else if laneID == enums.LaneRight {
		playerState.SetRightLaneCards(append(playerState.GetRightLaneCards(), cardInstance))
	}
}

func ReducePlayerHealth(playerState *models.PlayerMatchState, amount int) {
	var updatedHealth = playerState.GetHealth() - amount
	playerState.SetHealth(updatedHealth)

	var expectedRuneCount uint8 = uint8((updatedHealth - 1) / 5)
	var runeCount = max(0, min(expectedRuneCount, playerState.GetRunes()))
	playerState.SetRunes(runeCount)

	// TODO: trigger prophecies

	if updatedHealth <= 0 {
		// TODO: there'll be an exception with a Vivec card in play later
		match.EndMatch(playerState.MatchState, playerState.OpponentState.PlayerID)
	}
}

func ReduceCardHealth(playerState *models.PlayerMatchState, cardInstance *models.CardInstance, laneID enums.Lane, amount int) error {
	cardInstance.Health = cardInstance.Health - amount

	if cardInstance.Health <= 0 {
		DiscardCardFromLane(playerState, cardInstance, laneID)
	}

	playerState.SendEvent(enums.BackendEventCardInstancesChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentCardInstancesChanged)

	// to force lanes redraw
	playerState.SendEvent(enums.BackendEventLanesChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentLanesChanged)

	return nil
}

func StartTurn(playerState *models.PlayerMatchState) {
	playerState.SetMaxMana(playerState.GetMaxMana() + 1)
	playerState.SetMana(playerState.GetMaxMana())
	DrawCard(playerState)

	for _, card := range playerState.GetLeftLaneCards() {
		card.IsActive = true
	}
	for _, card := range playerState.GetRightLaneCards() {
		card.IsActive = true
	}
}
