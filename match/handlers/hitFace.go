package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/actions"
)

func HitFace(playerID int, cardInstanceID uuid.UUID) {
	_, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		return
	}

	cardInstance, _, _, err := match.GetCardInstanceFromLanes(playerState, cardInstanceID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		return
	}

	if !cardInstance.IsActive {
		fmt.Println(fmt.Errorf("[%d]: card with id '%s' is not active", playerID, cardInstanceID.String()))
		return
	}

	actions.ReducePlayerHealth(opponentState, cardInstance.Power)
	cardInstance.IsActive = false
}
