package senders

import (
	"fmt"

	"github.com/jental/freetesl-server/mappers"
	"github.com/jental/freetesl-server/models"
)

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
