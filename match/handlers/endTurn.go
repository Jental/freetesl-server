package handlers

import (
	"fmt"
	"time"

	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/actions"
	"github.com/jental/freetesl-server/match/senders"
)

func EndTurn(playerID int) {
	matchState, playerState, opponentState, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Println(err)
		return
	}

	actions.SwitchTurn(matchState)
	opponentState.MaxMana = opponentState.MaxMana + 1
	opponentState.Mana = opponentState.MaxMana
	actions.DrawCard(opponentState)
	actions.PlayRandomCards(opponentState)

	for _, card := range playerState.LeftLaneCards {
		card.IsActive = false
	}
	for _, card := range playerState.RightLaneCards {
		card.IsActive = false
	}

	senders.SendMatchStateToEveryone(matchState)

	time.Sleep(3 * time.Second)

	actions.SwitchTurn(matchState)
	actions.DrawCard(playerState)

	for _, card := range playerState.LeftLaneCards {
		card.IsActive = true
	}
	for _, card := range playerState.RightLaneCards {
		card.IsActive = true
	}

	playerState.MaxMana = playerState.MaxMana + 1
	playerState.Mana = playerState.MaxMana

	senders.SendMatchStateToEveryone(matchState)
}
