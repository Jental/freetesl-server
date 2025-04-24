package handlers

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/match/match"
	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/match/operations"
	"github.com/jental/freetesl-server/models/enums"
)

func MoveCardToLane(playerID int, cardInstanceID uuid.UUID, laneID enums.LanePosition, cardInstanceToReplaceID *uuid.UUID) {
	matchState, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Printf("[%d]: MoveCardToLane: %s", playerID, err)
		return
	}

	handleErr := func(err error) {
		fmt.Printf("[%d]: MoveCardToLane: %s", playerID, err)
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
	}

	cardInstance, _, exists := playerState.GetCardInstanceFromHand(cardInstanceID)
	if !exists {
		handleErr(fmt.Errorf("card instance with id '%s' is not present in a hand", cardInstanceID))
		return
	}

	creatureCardInstance, ok := cardInstance.(*models.CardInstanceCreature)
	if !ok {
		handleErr(errors.New("expected a creature card"))
		return
	}

	lane := playerState.GetLane(laneID)

	if cardInstanceToReplaceID != nil && lane.CountCardInstances() >= common.MAX_LANE_CARDS {
		// if num of cards lessaer than max we consider it an accidental replacement and will do regular move
		cardInstanceToReplace, _, exists := lane.GetCardInstance(*cardInstanceToReplaceID)
		if !exists {
			handleErr(fmt.Errorf("card instance with id '%s' is not present in a lane", cardInstanceToReplaceID))
			return
		}
		err = coreOperations.DiscardCardFromLane(playerState, cardInstanceToReplace, lane)
		if err != nil {
			handleErr(err)
			return
		}
	}

	err = operations.MoveCardFromHandToLane(playerState, matchState, creatureCardInstance, lane)
	if err != nil {
		handleErr(err)
		return
	}
}
