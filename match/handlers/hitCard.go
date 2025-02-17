package handlers

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/actions"
	"github.com/jental/freetesl-server/models/enums"
)

func HitCard(playerID int, cardInstanceID uuid.UUID, opponentCardInstanceID uuid.UUID) {
	_, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Println(err)
		return
	}

	cardInstance, laneID, _, err := match.GetCardInstanceFromLanes(playerState, cardInstanceID)
	if err != nil {
		fmt.Println(err)
		return
	}

	opponentCardInstance, opponentLaneID, _, err := match.GetCardInstanceFromLanes(opponentState, opponentCardInstanceID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !cardInstance.IsActive {
		fmt.Println(fmt.Errorf("card with id '%d' is not active", cardInstanceID))
		return
	}

	if laneID != opponentLaneID {
		fmt.Println(errors.New("cards are on different lanes"))
		return
	}

	err = actions.ReduceCardHealth(opponentState, opponentCardInstance, opponentLaneID, cardInstance.Power)
	if err != nil {
		fmt.Println(err)
		return
	}

	cardInstance.IsActive = false

	playerState.Events <- enums.BackendEventLanesChanged
	opponentState.Events <- enums.BackendEventLanesChanged
}
