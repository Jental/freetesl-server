package senders

import (
	"log"

	"github.com/jental/freetesl-server/match/mappers"
	"github.com/jental/freetesl-server/match/models"
)

func SendDeckStateToPlayer(playerState *models.PlayerMatchState, match *models.Match) error {
	if playerState.Connection == nil {
		return nil // Fake opponent has nil connection. TODO: the check should be removed
	}

	var dto, err = mappers.MapToDeckStateDTO(match, playerState.PlayerID)
	if err != nil {
		return err
	}
	var json = map[string]interface{}{
		"method": "deckStateUpdate",
		"body":   dto,
	}

	log.Printf("[%d]: sending: deckStateUpdate", playerState.PlayerID)

	// TODO: each active player should have two queues:
	// - of requests from client to be processed
	// - of messages from server
	//   ideally with some filtration to avoid sending multiple matchStates one after another
	err = sendJson(playerState, json)
	if err != nil {
		return err
	}

	log.Printf("[%d]: sent: deckStateUpdate", playerState.PlayerID)

	return nil
}
