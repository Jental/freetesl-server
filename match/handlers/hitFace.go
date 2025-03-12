package handlers

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/operations"
)

func HitFace(playerID int, cardInstanceID uuid.UUID) {
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

	err = operations.HitFace(playerState, opponentState, cardInstance, laneID)
	if err != nil {
		log.Println(err)
	}
}
