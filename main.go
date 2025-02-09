package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/dtos"
	"github.com/mitchellh/mapstructure"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", connectAndJoinMatch)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func connectAndJoinMatch(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade error:", err)
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
			var dto dtos.JoinRequestDTO
			mapstructure.Decode(body, &dto)
			go joinMatch(dto.PlayerID, common.Maybe[uuid.UUID]{HasValue: false}, c) // for now always joing to a new match. TODO: fix
		case "endTurn":
			var dto dtos.EndTurnRequestDTO
			mapstructure.Decode(body, &dto)
			go endTurn(dto.PlayerID)
		case "moveCardToLane":
			var dto dtos.MoveCardToLaneRequestDTO
			mapstructure.Decode(body, &dto)
			cardInstanceID, err := uuid.Parse(dto.CardInstanceID)
			if err != nil {
				log.Println(err)
				continue
			}
			go moveCardToLane(dto.PlayerID, cardInstanceID, dto.LaneID)
		}
	}
}
