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
	"github.com/jental/freetesl-server/dtos"
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

func joinMatch(playerID int, matchID common.Maybe[uuid.UUID], conn *websocket.Conn) error {
	var match models.Match
	if !matchID.HasValue {
		match = models.Match{
			Id:          uuid.New(),
			Connection0: conn,
			Connection1: nil,
			Player0ID:   common.Maybe[int]{Value: playerID, HasValue: true},
			Player1ID:   common.Maybe[int]{HasValue: false},
		}
		mutex.Lock()
		matches = make(map[uuid.UUID]models.Match)
		matches[match.Id] = match
		mutex.Unlock()
	} else {
		match, exist := getMatch(matchID.Value)
		if !exist {
			return errors.New("Match with given id does not exist")
		} else if match.Player0ID.HasValue && match.Player1ID.HasValue && match.Player0ID.Value != playerID && match.Player1ID.Value != playerID {
			return errors.New("Match is already started with different players")
		} else if match.Player0ID.HasValue && match.Player0ID.Value == playerID {
			match.Connection0 = conn // updating connection just in case
		} else if match.Player1ID.HasValue && match.Player1ID.Value == playerID {
			match.Connection1 = conn // updating connection just in case
		} else if !match.Player0ID.HasValue {
			match.Player0ID = common.Maybe[int]{Value: playerID, HasValue: true} // we are joining as first player
		} else if !match.Player1ID.HasValue {
			match.Player1ID = common.Maybe[int]{Value: playerID, HasValue: true} // we are joining as second player
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

func sendMatchState() error {
	var matchState = dtos.PlayerMatchStateDTO{Hand: hand, Deck: leftDeck, Health: health, Runes: runes, Mana: mana, MaxMana: maxMana, OwnTurn: ownTurn}
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

		var health = 22
		var runes uint8 = 4
		var mana = 2
		var maxMana = 8
		var ownTurn = true

		// var matchState = dtos.PlayerMatchStateDTO{Hand: hand, Deck: leftDeck, Health: health, Runes: runes, Mana: mana, MaxMana: maxMana, OwnTurn: ownTurn}
		// var json = map[string]interface{}{
		// 	"method": "matchStateUpdate",
		// 	"body":   matchState,
		// }

		// // TODO: each active player should have two queues:
		// // - of requests from client to be processed
		// // - of messages from server
		// //   ideally with some filtration to avoid sending multiple matchStates one after another
		// err2 := conn.WriteJSON(json)
		// if err2 != nil {
		// 	return err2
		// }

		time.Sleep(3 * time.Second)
	}

	return nil
}
