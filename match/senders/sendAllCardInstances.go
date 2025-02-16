package senders

import (
	"errors"
	"fmt"
	"slices"

	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/mappers"
	"github.com/jental/freetesl-server/models"
	"github.com/samber/lo"
)

func SendAllCardInstancesToEveryone(match *models.Match) {
	cards, err := getAllCardInstances(match)
	if err != nil {
		fmt.Println(err)
		return
	}

	var errChan = make(chan error)

	if match.Player0State.HasValue {
		go func() {
			errChan <- sendAllCardInstancesToPlayer(match.Player0State.Value, cards)
		}()
	}

	if match.Player1State.HasValue {
		go func() {
			errChan <- sendAllCardInstancesToPlayer(match.Player1State.Value, cards)
		}()
	}

	var aggErrors = []error{
		<-errChan,
		<-errChan,
	}

	var errorsPresent = slices.IndexFunc(aggErrors, func(err error) bool { return err != nil }) >= 0
	if errorsPresent {
		fmt.Println(errors.Join(aggErrors...))
	}
}

func sendAllCardInstancesToPlayer(playerState *models.PlayerMatchState2, cards []*models.CardInstance) error {
	if playerState.Connection == nil {
		return nil // Fake opponent has nil connection. TODO: the check should be removed
	}

	var dto = lo.Map(cards, func(card *models.CardInstance, _ int) dtos.CardInstanceDTO { return mappers.MapToCardInstanceDTO(card) })
	var json = map[string]interface{}{
		"method": "allCardInstances",
		"body":   dto,
	}

	// TODO: each active player should have two queues:
	// - of requests from client to be processed
	// - of messages from server
	//   ideally with some filtration to avoid sending multiple matchStates one after another
	err := playerState.Connection.WriteJSON(json)
	if err != nil {
		return err
	}

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

func getAllCardInstancesFromPlayer(playerState *models.PlayerMatchState2) []*models.CardInstance {
	var result []*models.CardInstance
	result = append(result, playerState.Deck...)
	result = append(result, playerState.Hand...)
	result = append(result, playerState.LeftLaneCards...)
	result = append(result, playerState.RightLaneCards...)
	result = append(result, playerState.DiscardPile...)
	return result
}
