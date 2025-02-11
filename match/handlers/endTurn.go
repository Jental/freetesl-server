package handlers

import (
	"fmt"
	"time"

	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/actions"
	"github.com/jental/freetesl-server/match/senders"
)

func EndTurn(playerID int) {
	matchState, playerState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Println(err)
		return
	}

	actions.SwitchTurn(matchState)
	senders.SendMatchStateToEveryone(matchState)

	time.Sleep(3 * time.Second)
	actions.SwitchTurn(matchState)
	actions.DrawCard(playerState)

	for _, card := range playerState.LeftLaneCards {
		card.IsActive = true
		match.CardInstanceLastEndTurned = card
	}

	playerState.MaxMana = playerState.MaxMana + 1
	playerState.Mana = playerState.MaxMana

	match.PlayerLastEndTurned = playerState

	senders.SendMatchStateToEveryone(matchState)
}
