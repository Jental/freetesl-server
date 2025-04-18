package senders

import (
	"log"

	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/db/queries"
	"github.com/jental/freetesl-server/mappers"
	"github.com/jental/freetesl-server/match/models"
)

func SendAllCardsToPlayer(playerState *models.PlayerMatchState) error {
	cards, err := queries.GetAllCards()
	if err != nil {
		return err
	}

	err = sendAllCardsToPlayer(playerState, cards)
	if err != nil {
		return err
	}

	return nil
}

func sendAllCardsToPlayer(playerState *models.PlayerMatchState, cards []*dbModels.Card) error {
	if playerState.Connection == nil {
		return nil // Fake opponent has nil connection. TODO: the check should be removed
	}

	var dto = mappers.MapToAllCardsDTO(cards)
	var json = map[string]interface{}{
		"method": "allCards",
		"body":   dto,
	}

	log.Printf("[%d]: sending: allCards", playerState.PlayerID)

	// TODO: each active player should have two queues:
	// - of requests from client to be processed
	// - of messages from server
	//   ideally with some filtration to avoid sending multiple matchStates one after another
	err := sendJson(playerState, json)
	if err != nil {
		return err
	}

	log.Printf("[%d]: sent: allCards", playerState.PlayerID)

	return nil
}
