package main

import (
	"errors"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/mappers"
	"github.com/jental/freetesl-server/models"
	"github.com/samber/lo"
)

var matches map[uuid.UUID]models.Match = nil
var once sync.Once
var mutex sync.Mutex

func getMatch(matchID uuid.UUID) (models.Match, bool) {
	once.Do(func() {
		var m = make(map[uuid.UUID]models.Match)
		matches = m
	})
	mutex.Lock()
	defer mutex.Unlock()

	match, exist := matches[matchID]
	return match, exist
}

func createInitialPlayerMatchState(playerID int, conn *websocket.Conn) (models.PlayerMatchState2, error) {
	decks, err := getDecks(playerID)
	if err != nil {
		return models.PlayerMatchState2{}, err
	}

	var deckInstance []models.CardInstance = lo.Shuffle(
		lo.FlatMap(
			slices.Collect(maps.Values(decks[0].Cards)),
			func(cardWithCount models.CardWithCount, _ int) []models.CardInstance {
				return lo.Times(cardWithCount.Count, func(_ int) models.CardInstance {
					return models.CardInstance{
						Card:           cardWithCount.Card,
						CardInstanceID: uuid.New(),
						Power:          cardWithCount.Card.Power,
						Health:         cardWithCount.Card.Health,
					}
				})
			}))

	var hand []models.CardInstance = deckInstance[:3]
	var leftDeck = deckInstance[3:]

	return models.PlayerMatchState2{
		PlayerID:   playerID,
		Connection: conn,
		Deck:       leftDeck,
		Hand:       hand,
		Health:     30,
		Runes:      5,
		Mana:       1,
		MaxMana:    1,
	}, nil
}

func joinMatch(playerID int, matchID common.Maybe[uuid.UUID], conn *websocket.Conn) error {
	var match models.Match
	if !matchID.HasValue {
		matchState, err := createInitialPlayerMatchState(playerID, conn)
		if err != nil {
			return err
		}

		match = models.Match{
			Id: uuid.New(),
			Player0State: common.Maybe[models.PlayerMatchState2]{
				HasValue: true,
				Value:    matchState,
			},
			Player1State: common.Maybe[models.PlayerMatchState2]{
				HasValue: false,
			},
		}

		mutex.Lock()
		matches = make(map[uuid.UUID]models.Match)
		matches[match.Id] = match
		mutex.Unlock()
	} else {
		match, exist := getMatch(matchID.Value)
		if !exist {
			return errors.New("Match with given id does not exist")
		} else if match.Player0State.HasValue && match.Player1State.HasValue && match.Player0State.Value.PlayerID != playerID && match.Player1State.Value.PlayerID != playerID {
			return errors.New("Match is already started with different players")
		} else if match.Player0State.HasValue && match.Player0State.Value.PlayerID == playerID {
			match.Player0State.Value.Connection = conn // updating connection just in case
		} else if match.Player1State.HasValue && match.Player1State.Value.PlayerID == playerID {
			match.Player1State.Value.Connection = conn // updating connection just in case
		} else if !match.Player0State.HasValue {
			// we are joining as first player
			matchState, err := createInitialPlayerMatchState(playerID, conn)
			if err != nil {
				return err
			}
			match.Player0State = common.Maybe[models.PlayerMatchState2]{
				HasValue: true,
				Value:    matchState,
			}
		} else if !match.Player1State.HasValue {
			// we are joining as second player
			matchState, err := createInitialPlayerMatchState(playerID, conn)
			if err != nil {
				return err
			}
			match.Player1State = common.Maybe[models.PlayerMatchState2]{
				HasValue: true,
				Value:    matchState,
			}
		}
	}

	go startTestCardDraw(match)

	return nil
}

func sendMatchStateToEveryone(match models.Match) error {
	if match.Player0State.HasValue {
		err := sendMatchStateToPlayer(match.Player0State.Value, match)
		if err != nil {
			return err
		}
	}

	if match.Player1State.HasValue {
		err := sendMatchStateToPlayer(match.Player1State.Value, match)
		if err != nil {
			return err
		}
	}

	return nil
}

func sendMatchStateToPlayer(playerState models.PlayerMatchState2, match models.Match) error {
	var ownTurn = playerState.PlayerID == match.PlayerWithTurnID
	var matchState = mappers.MapToPlayerMatchStateDTO(playerState, ownTurn)
	var json = map[string]interface{}{
		"method": "matchStateUpdate",
		"body":   matchState,
	}

	// TODO: each active player should have two queues:
	// - of requests from client to be processed
	// - of messages from server
	//   ideally with some filtration to avoid sending multiple matchStates one after another
	err := playerState.Connection.WriteJSON(json)
	if err != nil {
		return err
	}

	return nil
}

func startTestCardDraw(match models.Match) error {
	if !match.Player0State.HasValue {
		return nil
	}

	for {
		if len(match.Player0State.Value.Deck) == 0 {
			break
		}

		var drawnCard = match.Player0State.Value.Deck[0]
		match.Player0State.Value.Hand = append(match.Player0State.Value.Hand, drawnCard)
		match.Player0State.Value.Deck = match.Player0State.Value.Deck[1:]

		sendMatchStateToEveryone(match)

		time.Sleep(3 * time.Second)
	}

	return nil
}
