package handlers

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/operations"
	"github.com/jental/freetesl-server/models"
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
		errorResponseDTO := dtos.ErrorDTO{
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
	playerWithTurnID := selectRandomPlayer(playerID, opponentID)

	playerState, err := createInitialPlayerMatchState(playerID, playerWithTurnID == playerID, nil) // connections will be filled later when user establishes a connection
	if err != nil {
		return nil, err
	}

	opponentState, err := createInitialPlayerMatchState(opponentID, playerWithTurnID == opponentID, nil) // connections will be filled later when user establishes a connection
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

	services.SetPlayerState(playerID, enums.PlayerStateInMatch)
	services.SetPlayerState(opponentID, enums.PlayerStateInMatch)

	return &matchState.Id, nil
}

func createInitialPlayerMatchState(playerID int, hasFirstTurn bool, conn *websocket.Conn) (*models.PlayerMatchState, error) {
	decks, err := services.GetDecks(playerID)
	if err != nil {
		return nil, err
	}

	var deckInstance []*models.CardInstance = lo.Shuffle(
		lo.FlatMap(
			decks[0].Cards,
			func(cardWithCount *models.CardWithCount, _ int) []*models.CardInstance {
				return lo.Times(cardWithCount.Count, func(_ int) *models.CardInstance {
					return &models.CardInstance{
						Card:           cardWithCount.Card,
						CardInstanceID: uuid.New(),
						Power:          cardWithCount.Card.Power,
						Health:         cardWithCount.Card.Health,
						Cost:           cardWithCount.Card.Cost,
						Keywords:       cardWithCount.Card.Keywords,
						IsActive:       false,
						Effects:        make([]*models.Effect, 0),
					}
				})
			}))

	var hand []*models.CardInstance = deckInstance[:3]
	var leftDeck = deckInstance[3:]

	for _, card := range hand {
		card.IsActive = true
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
