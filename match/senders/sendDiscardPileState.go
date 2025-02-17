package senders

import (
	"fmt"

	"github.com/jental/freetesl-server/mappers"
	"github.com/jental/freetesl-server/models"
)

func SendDiscardPileToEveryone(match *models.Match) {
	if match.Player0State.HasValue {
		go sendDiscardPileStateToPlayerWithErrorHandling(match.Player0State.Value, match)
	}

	if match.Player1State.HasValue {
		go sendDiscardPileStateToPlayerWithErrorHandling(match.Player1State.Value, match)
	}
}

func sendDiscardPileStateToPlayerWithErrorHandling(playerState *models.PlayerMatchState2, match *models.Match) {
	var err = sendDiscardPileStateToPlayer(playerState, match)
	if err != nil {
		fmt.Println(err)
	}
}

func sendDiscardPileStateToPlayer(playerState *models.PlayerMatchState2, match *models.Match) error {
	if playerState.Connection == nil {
		return nil // Fake opponent has nil connection. TODO: the check should be removed
	}

	var dto, err = mappers.MapToDiscardPileStateDTO(match, playerState.PlayerID)
	if err != nil {
		return err
	}
	var json = map[string]interface{}{
		"method": "discardPileStateUpdate",
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

	return nil
}
