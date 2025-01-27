package main

import (
	"errors"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/samber/lo"
	"jental.name/tesl/dtos"
	"jental.name/tesl/models"
)

var matches map[uuid.UUID]Match = nil
var once sync.Once
var mutex sync.Mutex

type Match struct {
	id          uuid.UUID
	connection0 *websocket.Conn
	connection1 *websocket.Conn
	player0ID   Maybe[int]
	player1ID   Maybe[int]
}

func getMatch(matchID uuid.UUID) (Match, bool) {
	once.Do(func() {
		var m = make(map[uuid.UUID]Match)
		matches = m
	})
	mutex.Lock()
	defer mutex.Unlock()

	match, exist := matches[matchID]
	return match, exist
}

func joinMatch(playerID int, matchID Maybe[uuid.UUID], conn *websocket.Conn) error {
	var match Match
	if !matchID.HasValue {
		match = Match{
			id:          uuid.New(),
			connection0: conn,
			connection1: nil,
			player0ID:   Maybe[int]{Value: playerID, HasValue: true},
			player1ID:   Maybe[int]{HasValue: false},
		}
		mutex.Lock()
		matches = make(map[uuid.UUID]Match)
		matches[match.id] = match
		mutex.Unlock()
	} else {
		match, exist := getMatch(matchID.Value)
		if !exist {
			return errors.New("Match with given id does not exist")
		} else if match.player0ID.HasValue && match.player1ID.HasValue && match.player0ID.Value != playerID && match.player1ID.Value != playerID {
			return errors.New("Match is already started with different players")
		} else if match.player0ID.HasValue && match.player0ID.Value == playerID {
			match.connection0 = conn // updating connection just in case
		} else if match.player1ID.HasValue && match.player1ID.Value == playerID {
			match.connection1 = conn // updating connection just in case
		} else if !match.player0ID.HasValue {
			match.player0ID = Maybe[int]{Value: playerID, HasValue: true} // we are joining as first player
		} else if !match.player1ID.HasValue {
			match.player1ID = Maybe[int]{Value: playerID, HasValue: true} // we are joining as second player
		}
	}

	decks, err := getDecks(playerID)
	if err != nil {
		return err
	}

	var deckInstance []dtos.CardInstanceDTO = lo.Shuffle(
		lo.FlatMap(
			slices.Collect(maps.Values(decks[0].Cards)),
			func(cardWithCount models.CardWithCount, _ int) []dtos.CardInstanceDTO {
				return lo.Times(cardWithCount.Count, func(_ int) dtos.CardInstanceDTO {
					return dtos.CardInstanceDTO{CardID: cardWithCount.Card.ID, CardInstanceID: uuid.New()}
				})
			}))

	go startTestCardDraw(deckInstance, conn)

	return nil
}

func startTestCardDraw(deckInstance []dtos.CardInstanceDTO, conn *websocket.Conn) error {
	var hand []dtos.CardInstanceDTO = deckInstance[:3]
	var leftDeck = deckInstance[3:]
	for {
		if len(leftDeck) == 0 {
			break
		}

		var drawnCard = leftDeck[0]
		hand = append(hand, drawnCard)
		leftDeck = leftDeck[1:]

		var matchState = dtos.PlayerMatchStateDTO{Hand: hand, Deck: leftDeck}
		var json = map[string]interface{}{
			"method": "matchStateUpdate",
			"body":   matchState,
		}

		// TODO: each active player should have two queues:
		// - of requests from client to be processed
		// - of messages from server
		//   ideally with some filtration to avoid sending multiple matchStates one after another
		err2 := conn.WriteJSON(json)
		if err2 != nil {
			return err2
		}

		time.Sleep(3 * time.Second)
	}

	return nil
}
