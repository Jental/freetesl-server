package senders

import (
	"fmt"

	"github.com/jental/freetesl-server/mappers"
	"github.com/jental/freetesl-server/models"
)

func sendDeckToEveryone(match *models.Match) {
	if match.Player0State.HasValue {
		go sendDeckStateToPlayerWithErrorHandling(match.Player0State.Value, match)
	}

	if match.Player1State.HasValue {
		go sendDeckStateToPlayerWithErrorHandling(match.Player1State.Value, match)
	}
}

func sendDeckStateToPlayerWithErrorHandling(playerState *models.PlayerMatchState, match *models.Match) {
	var err = SendDeckStateToPlayer(playerState, match)
	if err != nil {
		fmt.Println(err)
	}
}

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

	// TODO: each active player should have two queues:
	// - of requests from client to be processed
	// - of messages from server
	//   ideally with some filtration to avoid sending multiple matchStates one after another
	err = playerState.Connection.WriteJSON(json)
	if err != nil {
		return err
	}

	fmt.Printf("sent: [%d]: deckStateUpdate\n", playerState.PlayerID)

	return nil
}
