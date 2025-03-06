package handlers

import (
	"fmt"

	"github.com/jental/freetesl-server/match"
)

func Concede(playerID int) {
	matchState, _, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		return
	}

	match.EndMatch(matchState, opponentState.PlayerID)
}
