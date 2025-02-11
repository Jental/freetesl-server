package match

import "github.com/jental/freetesl-server/models"

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
