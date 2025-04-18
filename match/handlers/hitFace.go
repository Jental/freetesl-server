package handlers

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/match/match"
	"github.com/jental/freetesl-server/match/operations"
)

func HitFace(playerID int, cardInstanceID uuid.UUID) {
	_, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		return
	}

	cardInstance, lane, _, exists := playerState.GetCardInstanceFromLanes(cardInstanceID)
	if !exists {
		fmt.Println(fmt.Errorf("[%d]: card instance with id '%s' is not present on lanes", playerState.PlayerID, cardInstanceID))
		return
	}

	err = operations.HitFace(playerState, opponentState, cardInstance, lane)
	if err != nil {
		log.Println(err)
	}
}
