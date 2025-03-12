package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/models/enums"
)

func MoveCardToLane(playerID int, cardInstanceID uuid.UUID, laneID enums.Lane) {
	_, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}

	cardInstance, idx, err := match.GetCardInstanceFromHand(playerState, cardInstanceID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}

	err = coreOperations.MoveCardToLane(playerState, cardInstance, idx, laneID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}
}
