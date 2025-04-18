package coreOperations

import (
	"errors"
	"slices"

	"github.com/jental/freetesl-server/match/models"
)

func DiscardCardFromHand(playerState *models.PlayerMatchState, cardInstance models.CardInstance) error {
	var hand = playerState.GetHand()
	var idx = slices.Index(hand, cardInstance)
	if idx < 0 {
		return errors.New("player does have the card in a hand")
	}
	playerState.SetHand(slices.Delete(hand, idx, idx+1))

	playerState.SetDiscardPile(append(playerState.GetDiscardPile(), cardInstance))

	return nil
}
