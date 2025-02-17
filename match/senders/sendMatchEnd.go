package senders

import (
	"fmt"

	"github.com/jental/freetesl-server/mappers"
	"github.com/jental/freetesl-server/models"
)

func sendMatchEndToEveryone(match *models.Match) {
	if match.Player0State.HasValue {
		go sendMatchEndToPlayerWithErrorHandling(match.Player0State.Value, match)
	}

	if match.Player1State.HasValue {
		go sendMatchEndToPlayerWithErrorHandling(match.Player1State.Value, match)
	}
}

func sendMatchEndToPlayerWithErrorHandling(playerState *models.PlayerMatchState, match *models.Match) {
	var err = SendMatchEndToPlayer(playerState, match)
	if err != nil {
		fmt.Println(err)
	}
}

func SendMatchEndToPlayer(playerState *models.PlayerMatchState, match *models.Match) error {
	if playerState.Connection == nil {
		return nil // Fake opponent has nil connection. TODO: the check should be removed
	}

	var dto = mappers.MapToMatchEndDTO(match, playerState.PlayerID)
	var json = map[string]interface{}{
		"method": "matchEnd",
		"body":   dto,
	}

	// TODO: each active player should have two queues:
	// - of requests from client to be processed
	// - of messages from server
	//   ideally with some filtration to avoid sending multiple matchStates one after another
	var err = playerState.Connection.WriteJSON(json)
	if err != nil {
		return err
	}

	fmt.Printf("sent: [%d]: matchEnd\n", playerState.PlayerID)

	return nil
}
