package actions

import (
	"github.com/jental/freetesl-server/match/senders"
	"github.com/jental/freetesl-server/models"
)

func DrawCard(playerState *models.PlayerMatchState2) {
	if len(playerState.Deck) == 0 {
		return // TODO: for now doing nothing, but later next rune should be broken
	}

	var drawnCard = playerState.Deck[0]
	drawnCard.IsActive = true
	playerState.Hand = append(playerState.Hand, drawnCard)
	playerState.Deck = playerState.Deck[1:]
}

func SwitchTurn(match *models.Match) {
	var isFirstPlayersTurn = match.Player0State.Value.PlayerID == match.PlayerWithTurnID
	if isFirstPlayersTurn {
		match.PlayerWithTurnID = match.Player1State.Value.PlayerID
	} else {
		match.PlayerWithTurnID = match.Player0State.Value.PlayerID
	}
}

func ReducePlayerHealth(playerState *models.PlayerMatchState2, matchState *models.Match, amount int) {
	playerState.Health = playerState.Health - amount

	var expectedRuneCount uint8 = uint8((playerState.Health - 1) / 5)
	var runeCount = max(0, min(expectedRuneCount, playerState.Runes))
	playerState.Runes = runeCount

	// TODO: trigger prophecies

	if playerState.Health <= 0 {
		senders.SendMatchEndToEveryone(matchState)
		// TODO: there'll be an exception with a Vivec card in play later
	}
}
