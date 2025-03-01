package handlers

import (
	"log"
	"net/http"

	"github.com/jental/freetesl-server/models/enums"
	"github.com/jental/freetesl-server/services"
)

func setPlayerState(w http.ResponseWriter, req *http.Request, state enums.PlayerState) {
	var playerID int = -1
	contextVal := req.Context().Value(enums.ContextKeyUserID)
	if contextVal == nil {
		log.Println("player id is not found in a context")
	} else {
		var ok bool = false
		playerID, ok = contextVal.(int)
		if !ok {
			log.Println("player id from a context has invalid type")
			playerID = -1 // to have a common error handling
		}
	}
	if playerID < 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	services.SetPlayerState(playerID, state)
	w.WriteHeader(http.StatusOK)
}

func StartLookingForOpponent(w http.ResponseWriter, req *http.Request) {
	setPlayerState(w, req, enums.PlayerStateLookingForOpponent)
}

func StopLookingForOpponent(w http.ResponseWriter, req *http.Request) {
	setPlayerState(w, req, enums.PlayerStateOnline)
}
