package handlers

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/actions"
	"github.com/jental/freetesl-server/match/senders"
)

func HitCard(playerID int, cardInstanceID uuid.UUID, opponentCardInstanceID uuid.UUID) {
	matchState, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Println(err)
		return
	}

	cardInstance, laneID, _, err := match.GetCardInstanceFromLanes(playerState, cardInstanceID)
	if err != nil {
		fmt.Println(err)
		return
	}

	opponentCardInstance, opponentLaneID, _, err := match.GetCardInstanceFromLanes(opponentState, opponentCardInstanceID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !cardInstance.IsActive {
		fmt.Println(fmt.Errorf("card with id '%d' is not active", cardInstanceID))
		return
	}

	if laneID != opponentLaneID {
		fmt.Println(errors.New("cards are on different lanes"))
		return
	}

	err = actions.ReduceCardHealth(opponentCardInstance, opponentLaneID, opponentState, matchState, cardInstance.Power)
	if err != nil {
		fmt.Println(err)
		return
	}

	cardInstance.IsActive = false

	senders.SendMatchStateToEveryone(matchState)
}
