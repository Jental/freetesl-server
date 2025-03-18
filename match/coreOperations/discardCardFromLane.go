package coreOperations

import (
	"errors"
	"slices"

	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

func DiscardCardFromLane(playerState *models.PlayerMatchState, cardInstance *models.CardInstance, laneID enums.Lane) error {
	if laneID == enums.LaneLeft {
		var idx = slices.Index(playerState.GetLeftLaneCards(), cardInstance)
		if idx < 0 {
			return errors.New("player does have the card in a left lane")
		}

		playerState.SetLeftLaneCards(slices.Delete(playerState.GetLeftLaneCards(), idx, idx+1))
		playerState.SetDiscardPile(append(playerState.GetDiscardPile(), cardInstance))

	} else {
		var idx = slices.Index(playerState.GetRightLaneCards(), cardInstance)
		if idx < 0 {
			return errors.New("player does have the card in a right lane")
		}

		playerState.SetRightLaneCards(slices.Delete(playerState.GetRightLaneCards(), idx, idx+1))
		playerState.SetDiscardPile(append(playerState.GetDiscardPile(), cardInstance))
	}

	return nil
}
