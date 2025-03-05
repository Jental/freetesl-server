package handlers

import (
	"encoding/json"
	"log"
	"maps"
	"math/rand"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/db"
	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/match"
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
		PlayerWithTurnID: selectRandomPlayer(playerID, opponentID),
		WinnerID:         -1,
	}

	match.AddOrRefreshMatch(&matchState)
	updateMatchPlayerFields(&matchState)

	services.SetPlayerState(playerID, enums.PlayerStateInMatch)
	services.SetPlayerState(opponentID, enums.PlayerStateInMatch)

	return &matchState.Id, nil
}

func createInitialPlayerMatchState(playerID int, conn *websocket.Conn) (models.PlayerMatchState, error) {
	decks, err := db.GetDecks(playerID)
	if err != nil {
		return models.PlayerMatchState{}, err
	}

	var deckInstance []*models.CardInstance = lo.Shuffle(
		lo.FlatMap(
			slices.Collect(maps.Values(decks[0].Cards)),
			func(cardWithCount dbModels.CardWithCount, _ int) []*models.CardInstance {
				return lo.Times(cardWithCount.Count, func(_ int) *models.CardInstance {
					return &models.CardInstance{
						Card:           cardWithCount.Card,
						CardInstanceID: uuid.New(),
						Power:          cardWithCount.Card.Power,
						Health:         cardWithCount.Card.Health,
						Cost:           cardWithCount.Card.Cost,
						Keywords:       cardWithCount.Card.Keywords,
						IsActive:       false,
					}
				})
			}))

	var hand []*models.CardInstance = deckInstance[:3]
	var leftDeck = deckInstance[3:]

	for _, card := range hand {
		card.IsActive = true
	}

	var playerState = models.NewPlayerMatchState(
		playerID,
		30,
		5,
		1,
		1,
		leftDeck,
		hand,
		[]*models.CardInstance{},
		[]*models.CardInstance{},
		[]*models.CardInstance{},
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
