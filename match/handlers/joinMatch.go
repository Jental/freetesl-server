package handlers

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"time"

	"math/rand"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/db"
	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/senders"
	"github.com/jental/freetesl-server/models"
	"github.com/samber/lo"
)

var rnd rand.Rand = *rand.New(rand.NewSource(time.Now().UnixNano()))

func createInitialPlayerMatchState(playerID int, conn *websocket.Conn) (models.PlayerMatchState2, error) {
	decks, err := db.GetDecks(playerID)
	if err != nil {
		return models.PlayerMatchState2{}, err
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
						IsActive:       false,
					}
				})
			}))

	var hand []*models.CardInstance = deckInstance[:3]
	var leftDeck = deckInstance[3:]

	for _, card := range hand {
		card.IsActive = true
	}

	return models.PlayerMatchState2{
		PlayerID:       playerID,
		Connection:     conn,
		Deck:           leftDeck,
		Hand:           hand,
		LeftLaneCards:  []*models.CardInstance{},
		RightLaneCards: []*models.CardInstance{},
		DiscardPile:    []*models.CardInstance{},
		Health:         30,
		Runes:          5,
		Mana:           1,
		MaxMana:        1,
	}, nil
}

func selectRandomPlayer(player0ID int, player1ID int) int {
	var idx = rnd.Intn(2)
	if idx == 0 {
		return player0ID
	} else {
		return player1ID
	}
}

func createNewMatchForPlayer(playerID int, conn *websocket.Conn) (*models.Match, error) {
	playerState, err := createInitialPlayerMatchState(playerID, conn)
	if err != nil {
		return nil, err
	}

	// TODO: for debug, shoud be removed later
	// this block auocreates opponent
	var opponentID = 2
	player1State, err := createInitialPlayerMatchState(opponentID, nil)
	if err != nil {
		return nil, err
	}

	var matchState = models.Match{
		Id: uuid.New(),
		Player0State: common.Maybe[models.PlayerMatchState2]{
			HasValue: true,
			Value:    &playerState,
		},
		Player1State: common.Maybe[models.PlayerMatchState2]{
			HasValue: true,
			Value:    &player1State,
		},
		// PlayerWithTurnID: common.NONE_PLAYER_ID,
		PlayerWithTurnID: playerID, // since we've created an opponent, match can be started. TODO: remove
		WinnerID:         -1,
	}

	match.AddOrRefreshMatch(&matchState)

	return &matchState, nil
}

func JoinMatch(playerID int, matchID common.Maybe[uuid.UUID], conn *websocket.Conn) {
	var matchState *models.Match
	var err error
	if !matchID.HasValue {
		matchState, err = createNewMatchForPlayer(playerID, conn)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		var exist bool
		matchState, exist = match.GetMatch(*matchID.Value)

		if !exist {
			fmt.Println(errors.New("match with given id does not exist"))
			return

		} else if matchState.Player0State.HasValue && matchState.Player1State.HasValue && matchState.Player0State.Value.PlayerID != playerID && matchState.Player1State.Value.PlayerID != playerID {
			fmt.Println(errors.New("match is already started with different players"))

		} else if matchState.Player0State.HasValue && matchState.Player0State.Value.PlayerID == playerID {
			matchState.Player0State.Value.Connection = conn // updating connection just in case

		} else if matchState.Player1State.HasValue && matchState.Player1State.Value.PlayerID == playerID {
			matchState.Player1State.Value.Connection = conn // updating connection just in case

		} else if !matchState.Player0State.HasValue {
			// we are joining as first player
			playerState, err := createInitialPlayerMatchState(playerID, conn)
			if err != nil {
				fmt.Println(err)
				return
			}
			matchState.Player0State = common.Maybe[models.PlayerMatchState2]{
				HasValue: true,
				Value:    &playerState,
			}

			if !matchState.Player1State.HasValue {
				matchState.PlayerWithTurnID = common.NONE_PLAYER_ID // match has only one player => noones turn
			} else {
				matchState.PlayerWithTurnID = selectRandomPlayer(matchState.Player0State.Value.PlayerID, matchState.Player1State.Value.PlayerID)
			}

			match.AddOrRefreshMatch(matchState)

		} else if !matchState.Player1State.HasValue {
			// we are joining as second player
			playerState, err := createInitialPlayerMatchState(playerID, conn)
			if err != nil {
				fmt.Println(err)
				return
			}
			matchState.Player1State = common.Maybe[models.PlayerMatchState2]{
				HasValue: true,
				Value:    &playerState,
			}

			if !matchState.Player0State.HasValue {
				matchState.PlayerWithTurnID = common.NONE_PLAYER_ID // match has only one player => noones turn
			} else {
				matchState.PlayerWithTurnID = selectRandomPlayer(matchState.Player0State.Value.PlayerID, matchState.Player1State.Value.PlayerID)
			}

			match.AddOrRefreshMatch(matchState)
		}
	}

	senders.SendAllCardsToEveryone(matchState)
	senders.SendMatchInformationToEveryone(matchState)
	senders.SendAllCardInstancesToEveryone(matchState)
	senders.SendMatchStateToEveryone(matchState)
	senders.SendDeckToEveryone(matchState)
	senders.SendDiscardPileToEveryone(matchState)
}
