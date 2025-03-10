package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jental/freetesl-server/mappers"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/models/enums"
	"github.com/jental/freetesl-server/services"
)

func GetCurrentPlayerInfo(w http.ResponseWriter, req *http.Request) {
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

	var player, err = services.GetPlayer(playerID)
	if err != nil {
		log.Println(err)
		return
	}

	if player.State == enums.PlayerStateInMatch {
		// we want actual information, not cached one
		_, _, _, err := match.GetCurrentMatchState(playerID)
		if err != nil {
			// match not found
			player.State = enums.PlayerStateOnline
		}
	}

	var responseDTO = mappers.MapToPlayerInformationDTO(player)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseDTO)
}
