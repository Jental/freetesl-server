package senders

import (
	"log"

	"github.com/jental/freetesl-server/match/mappers"
	"github.com/jental/freetesl-server/match/models"
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

	log.Printf("[%d]: sending: matchEnd", playerState.PlayerID)

	// TODO: each active player should have two queues:
	// - of requests from client to be processed
	// - of messages from server
	//   ideally with some filtration to avoid sending multiple matchStates one after another
	var err = sendJson(playerState, json)
	if err != nil {
		return err
	}

	log.Printf("[%d]: sent: matchEnd", playerState.PlayerID)

	return nil
}
