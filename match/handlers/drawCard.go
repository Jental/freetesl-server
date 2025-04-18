package handlers

import (
	"fmt"

	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/match/match"
	"github.com/jental/freetesl-server/models/enums"
)

func DrawCard(playerID int) {
	_, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}

	coreOperations.DrawCard(playerState)
}
