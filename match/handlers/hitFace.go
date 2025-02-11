package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/actions"
	"github.com/jental/freetesl-server/match/senders"
	"github.com/jental/freetesl-server/models"
)

func HitFace(playerID int, cardInstanceID uuid.UUID) {
	matchState, playerState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Println(err)
		return
	}

	var opponentState *models.PlayerMatchState2
	if matchState.Player0State.Value == playerState {
		opponentState = matchState.Player1State.Value
	} else if matchState.Player1State.Value == playerState {
		opponentState = matchState.Player0State.Value
	} else {
		fmt.Println(fmt.Errorf("player with id '%d' is not a part of a match", playerID))
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
