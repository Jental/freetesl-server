package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/senders"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
	"github.com/mitchellh/mapstructure"
)

var upgrader = websocket.Upgrader{} // use default options

func startListeningWebsocketMessages(playerState *models.PlayerMatchState) {
	playerID := playerState.PlayerID

	for {
		log.Printf("[%d]: reading next json\n", playerID)

		_, r, err := playerState.Connection.NextReader() // using NextReader instead NextJson for better error handling
		if err != nil {
			log.Printf("[%d]: websocket read error (or connection was closed on match end): '%s'", playerID, err)
			close(playerState.PartiallyParsedMessages)
			if playerState.Connection != nil {
				_ = playerState.Connection.Close() // just in case
			}
			return
		}

		var request map[string]interface{}
		err = json.NewDecoder(r).Decode(&request)
		if err == io.EOF {
			log.Printf("[%d]: websocket read error: one value is expected in the message", playerID)
			continue
		} else if err != nil {
			log.Printf("[%d]: websocket read error: Failed to parse json:'%s'", playerID, r)
			continue
		}
		log.Printf("[%d]: read json", playerID)

		method, exists := request["method"]
		if !exists {
			log.Printf("[%d]: websocket read error: unknown method", playerID)
			continue
		}

		body, exists := request["body"]
		if !exists {
			log.Printf("[%d]: websocket read error:  body is expected. method: %s", playerID, method)
		}
		log.Printf("[%d]: ws recv: %s\n", playerID, method)

		playerState.PartiallyParsedMessages <- models.PartiallyParsedMessage{
			Method: method.(string),
			Body:   body,
		}
	}
}

func startListeningPartiallyParsedMessages(playerState *models.PlayerMatchState) {
	playerID := playerState.PlayerID

	for {
		_, playerState, _, err := match.GetCurrentMatchState(playerID)
		if err != nil {
			log.Printf("[%d]: no active match for player", playerID)
			continue
		}

		log.Printf("[%d]: checking cancellation", playerID)
		select {
		case <-playerState.ConnectionCancellationChan:
			log.Printf("[%d]: cancel requested", playerID)
			go func() {
				time.Sleep(5 * time.Second) // to let some sends pass
				log.Printf("[%d]: cancelled", playerID)
				if playerState.Connection != nil {
					_ = playerState.Connection.Close()
				}
			}()
			return
		case message := <-playerState.PartiallyParsedMessages:
			err = processMessage(playerID, message)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func processMessage(playerID int, message models.PartiallyParsedMessage) error {
	log.Printf("[%d]: processing message: '%s'\n", playerID, message.Method)

	switch message.Method {
	case "endTurn":
		go EndTurn(playerID)
	case "moveCardToLane":
		var dto dtos.MoveCardToLaneRequestDTO
		mapstructure.Decode(message.Body, &dto)
		cardInstanceID, err := uuid.Parse(dto.CardInstanceID)
		if err != nil {
			return err
		}
		go MoveCardToLane(playerID, cardInstanceID, dto.LaneID)
	case "hitFace":
		var dto dtos.HitFaceDTO
		mapstructure.Decode(message.Body, &dto)
		cardInstanceID, err := uuid.Parse(dto.CardInstanceID)
		if err != nil {
			return err
		}
		go HitFace(playerID, cardInstanceID)
	case "hitCard":
		var dto dtos.HitCardDTO
		mapstructure.Decode(message.Body, &dto)
		cardInstanceID, err := uuid.Parse(dto.CardInstanceID)
		if err != nil {
			return err
		}
		opponentCardInstanceID, err := uuid.Parse(dto.OpponentCardInstanceID)
		if err != nil {
			return err
		}
		go HitCard(playerID, cardInstanceID, opponentCardInstanceID)
	}

	return nil
}

func ConnectAndJoinMatch(w http.ResponseWriter, req *http.Request) {
	contextVal := req.Context().Value(enums.ContextKeyUserID)
	if contextVal == nil {
		log.Println("player id is not found in a context")
		return
	}
	playerID, ok := contextVal.(int)
	if !ok {
		log.Printf("[%d]: player id from a context has invalid type\n", playerID)
		return
	}
	log.Printf("[%d]: connectAndJoinMatch", playerID)

	connection, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("[%d]: upgrade error: '%s'", playerID, err)
		return
	}

	matchState, playerState, _, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		log.Printf("[%d] get match error: '%s'", playerID, err)
		return
	}

	playerState.Connection = connection

	go startListeningWebsocketMessages(playerState)
	go startListeningPartiallyParsedMessages(playerState)
	go senders.StartListeningBackendEvents(playerState, matchState)
	go JoinMatch(playerState)
}
