package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/match/match"
	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/match/operations"
	"github.com/jental/freetesl-server/models/enums"
)

func ApplyCardToCard(playerID int, cardInstanceID uuid.UUID, targetCardInstanceID uuid.UUID) {
	matchState, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}

	cardInstance, _, exists := playerState.GetCardInstanceFromHand(cardInstanceID)
	if !exists {
		fmt.Println(fmt.Errorf("[%d]: ApplyCardToCard: card instance with id '%s' is not present in a hand", playerID, cardInstanceID))
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}

	isTargetCardFromOpponent := true
	targetCardInstance, targetLane, _, exists := opponentState.GetCardInstanceFromLanes(targetCardInstanceID)
	if !exists {
		isTargetCardFromOpponent = false
		targetCardInstance, targetLane, _, exists = playerState.GetCardInstanceFromLanes(targetCardInstanceID)
		if !exists {
			fmt.Printf("[%d]: ApplyCardToCard: card instance with id '%s' is not present", playerID, targetCardInstanceID)
			playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
			opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
			return
		}
	}

	err = nil
	switch castedCardInstance := cardInstance.(type) {
	case *models.CardInstanceAction:
		err = operations.PlayActionCard(playerState, opponentState, castedCardInstance, targetCardInstance, isTargetCardFromOpponent, targetLane)
	case *models.CardInstanceItem:
		err = operations.PlayItemCard(playerState, opponentState, matchState, castedCardInstance, targetCardInstance, isTargetCardFromOpponent)
	default:
		fmt.Printf("[%d]: ApplyCardToCard: Expected an action or item card", playerID)
		playerState.SendEvent(enums.BackendEventMatchStateRefresh)
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}
	if err != nil {
		fmt.Printf("[%d]: ApplyCardToCard: %s", playerID, err)
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}
}
