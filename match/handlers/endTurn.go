package handlers

import (
	"fmt"
	"time"

	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/senders"
)

func EndTurn(playerID int) {
	matchState, playerState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Println(err)
		return
	}

	match.SwitchTurn(matchState)
	senders.SendMatchStateToEveryone(matchState)

	time.Sleep(3 * time.Second)
	match.SwitchTurn(matchState)
	match.DrawCard(playerState)

	for _, card := range playerState.LeftLaneCards {
		card.IsActive = true
	}

	playerState.MaxMana = playerState.MaxMana + 1
	playerState.Mana = playerState.MaxMana

	senders.SendMatchStateToEveryone(matchState)
}
