package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/dtos"
	"github.com/mitchellh/mapstructure"
)

var upgrader = websocket.Upgrader{} // use default options

func ConnectAndJoinMatch(w http.ResponseWriter, req *http.Request) {
	contextVal := req.Context().Value("userID")
	if contextVal == nil {
		log.Println("player id is not found in a context")
		return
	}
	playerID, ok := contextVal.(int)
	if !ok {
		log.Println("player id from a context has invalid type")
		return
	}
	log.Printf("Player id: %d", playerID)

	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	defer c.Close()

	for {
		var request map[string]interface{}

		err := c.ReadJSON(&request)
		if err != nil {
			log.Println("websocket read error:", err)
			continue
		}

		method, exists := request["method"]
		if !exists {
			log.Println("websocket read error: unknown method")
			continue
		}

		body, exists := request["body"]
		if !exists {
			log.Printf("websocket read error:  body is expected. method: %s\n", method)
		}
		log.Printf("recv: %s\n", method)

		switch method {
		case "join":
			go JoinMatch(playerID, common.Maybe[uuid.UUID]{HasValue: false}, c) // for now always joing to a new match. TODO: fix
		case "endTurn":
			go EndTurn(playerID)
		case "moveCardToLane":
			var dto dtos.MoveCardToLaneRequestDTO
			mapstructure.Decode(body, &dto)
			cardInstanceID, err := uuid.Parse(dto.CardInstanceID)
			if err != nil {
				log.Println(err)
				continue
			}
			go MoveCardToLane(playerID, cardInstanceID, dto.LaneID)
		case "hitFace":
			var dto dtos.HitFaceDTO
			mapstructure.Decode(body, &dto)
			cardInstanceID, err := uuid.Parse(dto.CardInstanceID)
			if err != nil {
				log.Println(err)
				continue
			}
			go HitFace(playerID, cardInstanceID)
		case "hitCard":
			var dto dtos.HitCardDTO
			mapstructure.Decode(body, &dto)
			cardInstanceID, err := uuid.Parse(dto.CardInstanceID)
			if err != nil {
				log.Println(err)
				continue
			}
			opponentCardInstanceID, err := uuid.Parse(dto.OpponentCardInstanceID)
			if err != nil {
				log.Println(err)
				continue
			}

			go HitCard(playerID, cardInstanceID, opponentCardInstanceID)
		}
	}
}
