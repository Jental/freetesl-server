package handlers

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/operations"
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

	err = operations.HitCard(playerState, opponentState, cardInstance, opponentCardInstance, lane, opponentLane)
	if err != nil {
		log.Println(err)
	}
}
