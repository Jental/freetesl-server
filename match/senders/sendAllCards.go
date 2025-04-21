package senders

import (
	"log"

	"github.com/jental/freetesl-server/mappers"
	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/services"
)

func SendAllCardsToPlayer(playerState *models.PlayerMatchState) error {
	err := sendAllCardsToPlayer(playerState)
	if err != nil {
		return err
	}

	return nil
}

func sendAllCardsToPlayer(playerState *models.PlayerMatchState) error {
	if playerState.Connection == nil {
		return nil // Fake opponent has nil connection. TODO: the check should be removed
	}

	cards, err := services.GetAllCards()
	if err != nil {
		return err
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
	err = sendJson(playerState, json)
	if err != nil {
		return err
	}

	log.Printf("[%d]: sent: allCards", playerState.PlayerID)

	return nil
}
