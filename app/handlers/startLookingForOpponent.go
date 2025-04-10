package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/models/enums"
	"github.com/jental/freetesl-server/services"
)

func setPlayerState(w http.ResponseWriter, req *http.Request, state enums.PlayerState, selectedDeckID *int) {
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

	services.SetPlayerState(playerID, state, selectedDeckID)
	w.WriteHeader(http.StatusOK)
}

func StartLookingForOpponent(w http.ResponseWriter, req *http.Request) {
	var decoder = json.NewDecoder(req.Body)
	var dto dtos.StartLookingForOpponentDTO
	err := decoder.Decode(&dto)
	if err != nil {
		log.Panic(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	setPlayerState(w, req, enums.PlayerStateLookingForOpponent, &dto.DeckID)
}

func StopLookingForOpponent(w http.ResponseWriter, req *http.Request) {
	setPlayerState(w, req, enums.PlayerStateOnline, nil)
}
