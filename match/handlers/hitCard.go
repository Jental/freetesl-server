package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/actions"
	"github.com/jental/freetesl-server/models/enums"
)

func HitCard(playerID int, cardInstanceID uuid.UUID, opponentCardInstanceID uuid.UUID) {
	_, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		return
	}

	cardInstance, laneID, _, err := match.GetCardInstanceFromLanes(playerState, cardInstanceID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		return
	}

	opponentCardInstance, opponentLaneID, _, err := match.GetCardInstanceFromLanes(opponentState, opponentCardInstanceID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		return
	}

	if !cardInstance.IsActive {
		fmt.Println(fmt.Errorf("[%d]: card with id '%s' is not active", playerID, cardInstanceID.String()))
		return
	}

	if laneID != opponentLaneID {
		fmt.Println(fmt.Errorf("[%d]: cards are on different lanes", playerID))
		return
	}

	err = actions.ReduceCardHealth(opponentState, opponentCardInstance, opponentLaneID, cardInstance.Power)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		return
	}

	cardInstance.IsActive = false

	playerState.Events <- enums.BackendEventLanesChanged
	opponentState.Events <- enums.BackendEventLanesChanged
}
