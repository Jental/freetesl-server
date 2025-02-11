package handlers

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/senders"
	"github.com/jental/freetesl-server/models"
)

func MoveCardToLane(playerID int, cardInstanceID uuid.UUID, laneID byte) {
	matchState, playerState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Println(err)
		return
	}

	var idx = slices.IndexFunc(playerState.Hand, func(el *models.CardInstance) bool { return el.CardInstanceID == cardInstanceID })
	if idx < 0 {
		fmt.Println(fmt.Errorf("card instance with id '%s' is not present in a hand of a player '%d'", cardInstanceID, playerID))
		return
	}
	var cardInstance = playerState.Hand[idx]

	if cardInstance.Cost > playerState.Mana {
		fmt.Println(fmt.Errorf("not enough mana '%d' of '%d'", cardInstance.Cost, playerState.Mana))
		senders.SendMatchStateToEveryone(matchState)
		return
	}

	if laneID == common.LEFT_LANE_ID {
		if len(playerState.LeftLaneCards) >= common.MAX_LANE_CARDS {
			fmt.Println(fmt.Errorf("lane is already full"))
			return
		}
		playerState.LeftLaneCards = append(playerState.LeftLaneCards, cardInstance)
	} else if laneID == common.RIGHT_LANE_ID {
		// if len(playerState.RightLaneCards) >= common.MAX_LANE_CARDS {
		// 	fmt.Println(fmt.Errorf("lane is already full"))
		// 	return
		// }
		// playerState.RightLaneCards = append(playerState.LeftLaneCards, cardInstance)
	} else {
		fmt.Println(fmt.Errorf("invali lane id: %d", laneID))
		return
	}

	cardInstance.IsActive = false
	playerState.Hand = slices.Delete(playerState.Hand, idx, idx+1)
	playerState.Mana = playerState.Mana - cardInstance.Cost

	senders.SendMatchStateToEveryone(matchState)
}
