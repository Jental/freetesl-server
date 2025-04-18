package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jental/freetesl-server/common"
	commonDTOs "github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/match/dtos"
	"github.com/jental/freetesl-server/match/match"
	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/match/operations"
	"github.com/jental/freetesl-server/models/enums"
	"github.com/jental/freetesl-server/services"
	"github.com/samber/lo"
)

var rnd rand.Rand = *rand.New(rand.NewSource(time.Now().UnixNano()))

func MatchCreate(w http.ResponseWriter, req *http.Request) {
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

	_, _, _, err := match.GetCurrentMatchState(playerID)
	if err == nil {
		// player already have a match
		w.WriteHeader(http.StatusBadRequest)
		errorResponseDTO := commonDTOs.ErrorDTO{
			ErrorCode: int(enums.ErrorCodePlayerHasMatch),
			Message:   enums.ErrorCodeMessages[enums.ErrorCodePlayerHasMatch],
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponseDTO)
		return
	}

	var decoder = json.NewDecoder(req.Body)
	var dto dtos.MatchCreateDTO
	err = decoder.Decode(&dto)
	if err != nil {
		log.Panic(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	matchID, err := matchCreate(playerID, dto.OpponentID, dto.DeckID)
	if err != nil {
		log.Panic(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseDTO := commonDTOs.GuidIdDTO{
		ID: matchID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseDTO)
}

func matchCreate(playerID int, opponentID int, deckID int) (*uuid.UUID, error) {
	opponentRuntimeInfo, exists := services.GetPlayerRuntimeInfo(opponentID)
	if !exists || opponentRuntimeInfo.SelectedDeckID == nil {
		return nil, fmt.Errorf("[%d]: opponent selected deck id is not found", playerID)
	}

	playerWithTurnID := selectRandomPlayer(playerID, opponentID)

	playerState, err := createInitialPlayerMatchState(playerID, deckID, playerWithTurnID == playerID, nil) // connections will be filled later when user establishes a connection
	if err != nil {
		return nil, err
	}

	opponentState, err := createInitialPlayerMatchState(opponentID, *opponentRuntimeInfo.SelectedDeckID, playerWithTurnID == opponentID, nil) // connections will be filled later when user establishes a connection
	if err != nil {
		return nil, err
	}

	var matchState = models.Match{
		Id: uuid.New(),
		Player0State: common.Maybe[models.PlayerMatchState]{
			HasValue: true,
			Value:    playerState,
		},
		Player1State: common.Maybe[models.PlayerMatchState]{
			HasValue: true,
			Value:    opponentState,
		},
		TurnID:                0,
		PlayerWithTurnID:      playerWithTurnID,
		PlayerWithFirstTurnID: playerWithTurnID,
		WinnerID:              -1,
	}

	match.AddOrRefreshMatch(&matchState)
	updateMatchPlayerFields(&matchState)

	var playerWithTurn *models.PlayerMatchState
	if matchState.PlayerWithTurnID == playerID {
		playerWithTurn = playerState
	} else {
		playerWithTurn = opponentState
	}
	operations.StartTurn(playerWithTurn, &matchState)

	services.SetPlayerState(playerID, enums.PlayerStateInMatch, nil)
	services.SetPlayerState(opponentID, enums.PlayerStateInMatch, nil)

	return &matchState.Id, nil
}

func createInitialPlayerMatchState(playerID int, deckID int, hasFirstTurn bool, conn *websocket.Conn) (*models.PlayerMatchState, error) {
	deck, err := services.GetDeck(playerID, deckID)
	if err != nil {
		return nil, err
	}

	deckInstance := make([]models.CardInstance, 0)
	for _, cardWithCount := range deck.Cards {
		for i := 0; i < cardWithCount.Count; i++ {
			cardInstance, err := models.NewCardInstance(cardWithCount.Card)
			if err != nil {
				return nil, err
			}
			deckInstance = append(deckInstance, cardInstance)
		}
	}
	shuffledDeckInstance := lo.Shuffle(deckInstance)

	hand := shuffledDeckInstance[:3]
	leftDeck := shuffledDeckInstance[3:]

	for _, card := range hand {
		card.SetIsActive(true)
	}

	hasRing := !hasFirstTurn

	var playerState = models.NewPlayerMatchState(
		playerID,
		30,
		5,
		0,
		0,
		hasRing,
		3,
		leftDeck,
		hand,
		conn,
	)

	return playerState, nil
}

func selectRandomPlayer(player0ID int, player1ID int) int {
	var idx = rnd.Intn(2)
	if idx == 0 {
		return player0ID
	} else {
		return player1ID
	}
}

func updateMatchPlayerFields(matchState *models.Match) {
	if matchState.Player0State.HasValue {
		matchState.Player0State.Value.MatchState = matchState
		if matchState.Player1State.HasValue {
			matchState.Player0State.Value.OpponentState = matchState.Player1State.Value
		}
	}

	if matchState.Player1State.HasValue {
		matchState.Player1State.Value.MatchState = matchState
		if matchState.Player0State.HasValue {
			matchState.Player1State.Value.OpponentState = matchState.Player0State.Value
		}
	}
}
