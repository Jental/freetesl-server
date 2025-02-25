package senders

import (
	"errors"
	"fmt"
	"slices"

	"github.com/jental/freetesl-server/db"
	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/mappers"
	"github.com/jental/freetesl-server/models"
)

func SendAllCardsToEveryone(match *models.Match) {
	cards, err := db.GetAllCards()
	if err != nil {
		fmt.Println(err)
		return
	}

	var errChan = make(chan error)

	if match.Player0State.HasValue {
		go func() {
			errChan <- sendAllCardsToPlayer(match.Player0State.Value, cards)
		}()
	}

	if match.Player1State.HasValue {
		go func() {
			errChan <- sendAllCardsToPlayer(match.Player1State.Value, cards)
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

func SendAllCardsToPlayer(playerState *models.PlayerMatchState) {
	cards, err := db.GetAllCards()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = sendAllCardsToPlayer(playerState, cards)
	if err != nil {
		fmt.Println(err)
	}
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

	// TODO: each active player should have two queues:
	// - of requests from client to be processed
	// - of messages from server
	//   ideally with some filtration to avoid sending multiple matchStates one after another
	err := playerState.Connection.WriteJSON(json)
	if err != nil {
		return err
	}

	fmt.Printf("sent: [%d]: allCards\n", playerState.PlayerID)

	return nil
}
