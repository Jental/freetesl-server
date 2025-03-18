package handlers

import (
	"fmt"

	"github.com/google/uuid"
	dbEnums "github.com/jental/freetesl-server/db/enums"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/operations"
	"github.com/jental/freetesl-server/models/enums"
)

func ApplyActionToCard(playerID int, cardInstanceID uuid.UUID, targetCardInstanceID uuid.UUID) {
	_, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}

	cardInstance, _, err := match.GetCardInstanceFromHand(playerState, cardInstanceID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}

	isTargetCardFromOpponent := true
	targetCardInstance, targetLaneID, _, err := match.GetCardInstanceFromLanes(opponentState, targetCardInstanceID)
	if err != nil {
		isTargetCardFromOpponent = false
		targetCardInstance, targetLaneID, _, err = match.GetCardInstanceFromLanes(playerState, targetCardInstanceID)
		if err != nil {
			fmt.Printf("[%d]: %s", playerID, err)
			playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
			opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
			return
		}
	}

	if cardInstance.Card.Type != dbEnums.CardTypeAction {
		fmt.Printf("[%d]: %s", playerID, err)
		playerState.SendEvent(enums.BackendEventMatchStateRefresh)
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}

	err = operations.PlayActionCard(playerState, opponentState, cardInstance, targetCardInstance, isTargetCardFromOpponent, &targetLaneID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}
}
