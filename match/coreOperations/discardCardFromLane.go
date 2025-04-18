package coreOperations

import (
	"github.com/jental/freetesl-server/match/models"
)

func DiscardCardFromLane(playerState *models.PlayerMatchState, cardInstance *models.CardInstanceCreature, lane *models.Lane) error {
	err := lane.RemoveCardInstance(cardInstance)
	if err != nil {
		return err
	}
	playerState.SetDiscardPile(append(playerState.GetDiscardPile(), cardInstance))

	return nil
}
