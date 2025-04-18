package match

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/models/enums"
)

// Some cases better to test:
// + create a match and play it till one of players win
// + both players joined match (hands are visible), one concedes
// + only one player joined match (second haven't had game window with a match active) and this player concedes; after that second player activates his window
//     now it shows "Failed to connect" for second user - but for now it's fine
// + only one player joined match (second haven't had game window with a match active) and this player does some action (play a card / switch turn). Then second joins.
// + reconnect to unfinished match
// + connect both players to match. One one players turn disconnect other one. Play a card / switch turn
// - create new match with the same players after concede

var MatchMessageHandlerFn func(playerID int, message models.PartiallyParsedMessage) error = nil
var BackendEventHandlerFn func(playerState *models.PlayerMatchState, event enums.BackendEventType) error = nil

func JoinMatch(playerState *models.PlayerMatchState) {
	playerState.SendEvent(enums.BackendEventMatchStart)
	playerState.SendEvent(enums.BackendEventCardInstancesChanged)
	playerState.SendEvent(enums.BackendEventOpponentCardInstancesChanged)
	playerState.SendEvent(enums.BackendEventHandChanged)
	playerState.SendEvent(enums.BackendEventOpponentHandChanged)
	playerState.SendEvent(enums.BackendEventDeckChanged)
	playerState.SendEvent(enums.BackendEventOpponentDeckChanged)
	playerState.SendEvent(enums.BackendEventDiscardPileChanged)
	playerState.SendEvent(enums.BackendEventOpponentDiscardPileChanged)
}

func DisconnectFromMatch(playerState *models.PlayerMatchState) {
	// TODO: check: it may happen, that in time we are trying to close a connection it's used (or maybe even recreated after relogin)
	log.Printf("[%d]: closing websocket connection", playerState.PlayerID)
	if playerState.Connection == nil {
		log.Printf("[%d]: websocket connection in not present (already closed or not established yet)", playerState.PlayerID)
		return
	}
	err := playerState.Connection.Close()
	if err != nil {
		log.Printf("[%d]: error during websocket connection close: '%s'", playerState.PlayerID, err)
	}
}

func InitListenersAfterConnectionEstablished(playerState *models.PlayerMatchState) {
	log.Printf("[%d]: after connection init", playerState.PlayerID)

	clearUnhandledBackendEvents(playerState)

	playerState.PartiallyParsedMessages = make(chan models.PartiallyParsedMessage)
	playerState.Events = make(chan enums.BackendEventType, 10)

	go startListeningBackendEvents(playerState)
	go startListeningWebsocketMessages(playerState)
	go startListeningPartiallyParsedMessages(playerState)
}

func cleanupAfterConnectionClose(playerState *models.PlayerMatchState) {
	defer func() { // try/catch analog - to recover from second close of channel
		if r := recover(); r != nil {
			log.Printf("[%d]: cleaning up after disconnect: recovered from error: %s", playerState.PlayerID, r)
		}
	}()

	log.Printf("[%d]: cleaning up after disconnect", playerState.PlayerID)
	playerState.Connection = nil
	close(playerState.PartiallyParsedMessages)
	close(playerState.Events)
	close(playerState.WaitingForUserActionChan)
	log.Printf("[%d]: cleaning up after disconnect - done", playerState.PlayerID)
}

func startListeningWebsocketMessages(playerState *models.PlayerMatchState) {
	playerID := playerState.PlayerID

	for {
		log.Printf("[%d]: ws: reading next json", playerID)

		_, r, err := playerState.Connection.NextReader() // using NextReader instead NextJson for better error handling
		if err != nil {
			log.Printf("[%d]: ws read error / or connection was closed by client / or connection was closed on match end: '%s'", playerID, err)
			cleanupAfterConnectionClose(playerState)
			return
		}

		var request map[string]interface{}
		err = json.NewDecoder(r).Decode(&request)
		if err == io.EOF {
			log.Printf("[%d]: ws read error: one value is expected in the message", playerID)
			continue
		} else if err != nil {
			log.Printf("[%d]: ws read error: Failed to parse json:'%s'", playerID, r)
			continue
		}
		log.Printf("[%d]: read json", playerID)

		method, exists := request["method"]
		if !exists {
			log.Printf("[%d]: ws read error: unknown method", playerID)
			continue
		}

		body, exists := request["body"]
		if !exists {
			log.Printf("[%d]: websocket read error:  body is expected. method: %s", playerID, method)
		}
		log.Printf("[%d]: ws recv: '%s'", playerID, method)

		playerState.PartiallyParsedMessages <- models.PartiallyParsedMessage{
			Method: method.(string),
			Body:   body,
		}
	}
}

func startListeningPartiallyParsedMessages(playerState *models.PlayerMatchState) {
	playerID := playerState.PlayerID

	if MatchMessageHandlerFn == nil {
		err := fmt.Errorf("[%d]: MatchMessageHandlerFn is not set", playerID)
		log.Panic(err)
		panic(err)
	}

	for message := range playerState.PartiallyParsedMessages {
		_, _, _, err := GetCurrentMatchState(playerID)
		if err != nil {
			log.Printf("[%d]: no active match for player. closing ws connection", playerID)
			if playerState.Connection != nil {
				_ = playerState.Connection.Close()
			}
			return
		}

		err = MatchMessageHandlerFn(playerID, message)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	log.Printf("[%d]: messages channel is closed", playerID)
}

func clearUnhandledBackendEvents(playerState *models.PlayerMatchState) {
	var count = 0
	for len(playerState.Events) > 0 {
		<-playerState.Events
		count = count + 1
	}
	log.Printf("[%d]: cleared %d unhandled events", playerState.PlayerID, count)
}

func startListeningBackendEvents(playerState *models.PlayerMatchState) {
	if BackendEventHandlerFn == nil {
		err := fmt.Errorf("[%d]: BackendEventHandlerFn is not set", playerState.PlayerID)
		log.Panic(err)
		panic(err)
	}

	for event := range playerState.Events {
		var err = BackendEventHandlerFn(playerState, event)

		if err != nil {
			log.Printf("[%d]: sending error: '%s'", playerState.PlayerID, err)
		}
	}
}
