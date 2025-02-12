package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/actions"
	"github.com/jental/freetesl-server/match/senders"
)

func MoveCardToLane(playerID int, cardInstanceID uuid.UUID, laneID byte) {
	matchState, playerState, _, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Println(err)
		senders.SendMatchStateToEveryone(matchState)
		return
	}

	cardInstance, idx, err := match.GetCardInstanceFromHand(playerState, cardInstanceID)
	if err != nil {
		fmt.Println(err)
		senders.SendMatchStateToEveryone(matchState)
		return
	}

	err = actions.MoveCardToLane(playerState, cardInstance, idx, laneID)
	if err != nil {
		fmt.Println(err)
		senders.SendMatchStateToEveryone(matchState)
		return
	}

	senders.SendMatchStateToEveryone(matchState)
}
