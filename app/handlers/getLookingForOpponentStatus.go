package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/match"
)

func GetLookingForOpponentStatus(w http.ResponseWriter, req *http.Request) {
	var playerID int = -1
	contextVal := req.Context().Value("userID")
	var err error
	if contextVal == nil {
		err = errors.New("player id is not found in a context")
	} else {
		var ok bool = false
		playerID, ok = contextVal.(int)
		if !ok {
			err = errors.New("player id from a context has invalid type")
			playerID = -1 // to have a common error handling
		}
	}
	if playerID < 0 {
		w.WriteHeader(http.StatusUnauthorized)
		log.Panic(err)
		return
	}

	matchState, _, _, err := match.GetCurrentMatchState(playerID)
	var responseDTO dtos.GuidIdDTO
	if err != nil {
		log.Printf("GetLookingForOpponentStatus: '%s'", err)
		responseDTO = dtos.GuidIdDTO{
			ID: nil,
		}
	} else {
		responseDTO = dtos.GuidIdDTO{
			ID: &matchState.Id,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseDTO)
}
