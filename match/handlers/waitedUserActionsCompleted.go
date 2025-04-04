package handlers

import (
	"fmt"

	"github.com/jental/freetesl-server/match"
)

func WaitedUserActionsCompleted(playerID int) {
	_, playerState, _, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		return
	}

	playerState.WaitingForUserActionChan <- struct{}{}
}
