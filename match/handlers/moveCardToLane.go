package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/operations"
	"github.com/jental/freetesl-server/models/enums"
)

func MoveCardToLane(playerID int, cardInstanceID uuid.UUID, laneID enums.LanePosition) {
	matchState, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}

	cardInstance, idx, exists := playerState.GetCardInstanceFromHand(cardInstanceID)
	if !exists {
		fmt.Println(fmt.Errorf("[%d]: card instance with id '%s' is not present in a hand", playerID, cardInstanceID))
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}

	lane := playerState.GetLane(laneID)

	err = operations.MoveCardFromHandToLane(playerState, matchState, cardInstance, idx, lane)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}
}
