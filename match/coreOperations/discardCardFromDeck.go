package coreOperations

import (
	"github.com/jental/freetesl-server/match/models"
)

func DiscardCardFromDeck(playerState *models.PlayerMatchState) {
	var deck = playerState.GetDeck()
	if len(deck) == 0 {
		return // TODO: for now doing nothing, but later next rune should be broken
	}

	var drawnCard = deck[0]
	playerState.SetDeck(deck[1:])

	playerState.SetDiscardPile(append(playerState.GetDiscardPile(), drawnCard))
}
