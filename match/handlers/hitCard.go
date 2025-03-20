package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/models/enums"
)

func HitCard(playerID int, cardInstanceID uuid.UUID, opponentCardInstanceID uuid.UUID) {
	_, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Println(fmt.Errorf("[%d]: %s", playerID, err))
		return
	}

	cardInstance, lane, _, exists := playerState.GetCardInstanceFromLanes(cardInstanceID)
	if !exists {
		fmt.Println(fmt.Errorf("[%d]: card instance with id '%s' is not present on lanes", playerState.PlayerID, cardInstanceID))
		return
	}

	opponentCardInstance, opponentLane, _, exists := opponentState.GetCardInstanceFromLanes(opponentCardInstanceID)
	if !exists {
		fmt.Println(fmt.Errorf("[%d]: card instance with id '%s' is not present on opponent lanes", playerState.PlayerID, opponentCardInstanceID))
		playerState.SendEvent(enums.BackendEventLanesChanged) // to reset FE state
		opponentState.SendEvent(enums.BackendEventOpponentLanesChanged)
		return
	}

	if !cardInstance.IsActive {
		fmt.Println(fmt.Errorf("[%d]: card with id '%s' is not active", playerID, cardInstanceID.String()))
		playerState.SendEvent(enums.BackendEventLanesChanged) // to reset FE state
		opponentState.SendEvent(enums.BackendEventOpponentLanesChanged)
		return
	}

	if lane.Position != opponentLane.Position {
		fmt.Println(fmt.Errorf("[%d]: cards are on different lanes", playerID))
		playerState.SendEvent(enums.BackendEventLanesChanged) // to reset FE state
		opponentState.SendEvent(enums.BackendEventOpponentLanesChanged)
		return
	}

	coreOperations.ReduceCardHealth(opponentState, opponentCardInstance, opponentLane, cardInstance.Power)
	coreOperations.ReduceCardHealth(playerState, cardInstance, lane, opponentCardInstance.Power)

	cardInstance.IsActive = false

	playerState.SendEvent(enums.BackendEventLanesChanged)
	playerState.SendEvent(enums.BackendEventOpponentLanesChanged)
	opponentState.SendEvent(enums.BackendEventLanesChanged)
	opponentState.SendEvent(enums.BackendEventOpponentLanesChanged)
}
