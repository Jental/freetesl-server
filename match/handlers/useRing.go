package handlers

import (
	"fmt"

	"github.com/jental/freetesl-server/match/match"
	"github.com/jental/freetesl-server/models/enums"
)

func UseRing(playerID int) {
	_, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Printf("[%d]: %s", playerID, err)
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}

	if !playerState.HasRing() {
		fmt.Println(fmt.Errorf("[%d]: UseRing: player has no ring", playerID))
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}

	if !playerState.IsRingActive() {
		fmt.Println(fmt.Errorf("[%d]: UseRing: ring is inactive", playerID))
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}

	gemCount := playerState.GetRingGemCount()
	if gemCount == 0 {
		fmt.Println(fmt.Errorf("[%d]: UseRing: no gems left", playerID))
		playerState.SendEvent(enums.BackendEventMatchStateRefresh) // on UI card may be already moved. In this case we need to send match state to FE to reset UI state
		opponentState.SendEvent(enums.BackendEventOpponentMatchStateRefresh)
		return
	}

	playerState.SetMana(playerState.GetMana() + 1)
	playerState.SetRingGemCount(gemCount - 1)
	playerState.SetRingActivity(false)
}
