package senders

import (
	"fmt"
	"log"

	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/mappers"
	"github.com/jental/freetesl-server/models"
	"github.com/samber/lo"
)

func SendAllCardInstancesToPlayer(playerState *models.PlayerMatchState, matchState *models.Match) error {
	cards, err := getAllCardInstances(matchState)
	if err != nil {
		return err
	}
	// TODO: SendAllCardInstancesToPlayer will be called twice - for player and for opponent => getAllCardInstances will be called twice too, which is not good
	//       maybe, we can pass instances with events

	return sendAllCardInstancesToPlayer(playerState, cards)
}

func sendAllCardInstancesToPlayer(playerState *models.PlayerMatchState, cards []*models.CardInstance) error {
	if playerState.Connection == nil {
		return nil // Fake opponent has nil connection. TODO: the check should be removed
	}

	var dto = lo.Map(cards, func(card *models.CardInstance, _ int) dtos.CardInstanceDTO { return mappers.MapToCardInstanceDTO(card) })
	var json = map[string]interface{}{
		"method": "allCardInstances",
		"body":   dto,
	}

	log.Printf("[%d]: sending: allCardInstances", playerState.PlayerID)

	// TODO: each active player should have two queues:
	// - of requests from client to be processed
	// - of messages from server
	//   ideally with some filtration to avoid sending multiple matchStates one after another
	err := sendJson(playerState, json)
	if err != nil {
		return err
	}

	log.Printf("[%d]: sent: allCardInstances", playerState.PlayerID)

	return nil
}

func getAllCardInstances(matchState *models.Match) ([]*models.CardInstance, error) {
	if !matchState.Player0State.HasValue || !matchState.Player1State.HasValue {
		return nil, fmt.Errorf("match is not started yet")
	}

	var player0CardInstances = getAllCardInstancesFromPlayer(matchState.Player0State.Value)
	var player1CardInstances = getAllCardInstancesFromPlayer(matchState.Player1State.Value)

	return append(player0CardInstances, player1CardInstances...), nil
}

func getAllCardInstancesFromPlayer(playerState *models.PlayerMatchState) []*models.CardInstance {
	var result []*models.CardInstance
	result = append(result, playerState.GetDeck()...)
	result = append(result, playerState.GetHand()...)
	result = append(result, playerState.GetLeftLaneCards()...)
	result = append(result, playerState.GetRightLaneCards()...)
	result = append(result, playerState.GetDiscardPile()...)
	return result
}
