package senders

import (
	"fmt"

	"github.com/jental/freetesl-server/db"
	"github.com/jental/freetesl-server/mappers"
	"github.com/jental/freetesl-server/models"
)

func SendAllCardsToEveryone(match *models.Match) {
	if match.Player0State.HasValue {
		go sendAllCardsToPlayerrWithErrorHandling(match.Player0State.Value)
	}

	if match.Player1State.HasValue {
		go sendAllCardsToPlayerrWithErrorHandling(match.Player1State.Value)
	}
}

func sendAllCardsToPlayerrWithErrorHandling(playerState *models.PlayerMatchState2) {
	var err = sendAllCardsToPlayer(playerState)
	if err != nil {
		fmt.Println(err)
	}
}

func sendAllCardsToPlayer(playerState *models.PlayerMatchState2) error {
	if playerState.Connection == nil {
		return nil // Fake opponent has nil connection. TODO: the check should be removed
	}

	cards, err := db.GetAllCards()
	if err != nil {
		return err
	}
	var dto = mappers.MapToAllCardsDTO(cards)
	var json = map[string]interface{}{
		"method": "allCards",
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
