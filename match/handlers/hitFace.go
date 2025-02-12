package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/actions"
	"github.com/jental/freetesl-server/match/senders"
)

func HitFace(playerID int, cardInstanceID uuid.UUID) {
	matchState, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Println(err)
		return
	}

	cardInstance, _, _, err := match.GetCardInstanceFromLanes(playerState, cardInstanceID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !cardInstance.IsActive {
		fmt.Println(fmt.Errorf("card with id '%d' is not active", cardInstanceID))
		return
	}

	actions.ReducePlayerHealth(opponentState, matchState, cardInstance.Power)
	cardInstance.IsActive = false

	senders.SendMatchStateToEveryone(matchState)
}
