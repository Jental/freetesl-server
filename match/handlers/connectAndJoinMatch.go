package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/models/enums"
)

var upgrader = websocket.Upgrader{} // use default options

func ConnectAndJoinMatch(w http.ResponseWriter, req *http.Request) {
	contextVal := req.Context().Value(enums.ContextKeyUserID)
	if contextVal == nil {
		log.Println("player id is not found in a context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	playerID, ok := contextVal.(int)
	if !ok {
		log.Printf("[%d]: player id from a context has invalid type\n", playerID)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Printf("[%d]: connectAndJoinMatch", playerID)

	_, playerState, _, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		log.Printf("[%d] get match error: '%s'", playerID, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	connection, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("[%d]: upgrade error: '%s'", playerID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	playerState.Connection = connection

	match.InitListenersAfterConnectionEstablished(playerState)
	go match.JoinMatch(playerState)

	w.WriteHeader(http.StatusOK)
}
