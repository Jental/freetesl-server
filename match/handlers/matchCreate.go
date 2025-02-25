package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
	"github.com/jental/freetesl-server/services"
)

func MatchCreate(w http.ResponseWriter, req *http.Request) {
	var playerID int = -1
	contextVal := req.Context().Value("userID")
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

	var decoder = json.NewDecoder(req.Body)
	var dto dtos.MatchCreateDTO
	err := decoder.Decode(&dto)
	if err != nil {
		log.Panic(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	matchID, err := matchCreate(playerID, dto.OpponentID)
	if err != nil {
		log.Panic(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseDTO := dtos.GuidIdDTO{
		ID: matchID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseDTO)
}

func matchCreate(playerID int, opponentID int) (*uuid.UUID, error) {
	playerState, err := createInitialPlayerMatchState(playerID, nil) // connections will be filled later when user establishes a connection
	if err != nil {
		return nil, err
	}

	opponentState, err := createInitialPlayerMatchState(opponentID, nil) // connections will be filled later when user establishes a connection
	if err != nil {
		return nil, err
	}

	var matchState = models.Match{
		Id: uuid.New(),
		Player0State: common.Maybe[models.PlayerMatchState]{
			HasValue: true,
			Value:    &playerState,
		},
		Player1State: common.Maybe[models.PlayerMatchState]{
			HasValue: true,
			Value:    &opponentState,
		},
		PlayerWithTurnID: playerID, // TODO: random
		WinnerID:         -1,
	}

	match.AddOrRefreshMatch(&matchState)
	updateMatchPlayerFields(&matchState)

	services.SetPlayerState(playerID, enums.PlayerStateInMatch)
	services.SetPlayerState(opponentID, enums.PlayerStateInMatch)

	return &matchState.Id, nil
}
