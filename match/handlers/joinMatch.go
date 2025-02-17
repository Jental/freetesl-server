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
	"github.com/jental/freetesl-server/models/enums"
	"github.com/samber/lo"
)

var rnd rand.Rand = *rand.New(rand.NewSource(time.Now().UnixNano()))

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
		Player0State: common.Maybe[models.PlayerMatchState]{
			HasValue: true,
			Value:    &playerState,
		},
		Player1State: common.Maybe[models.PlayerMatchState]{
			HasValue: true,
			Value:    &player1State,
		},
		// PlayerWithTurnID: common.NONE_PLAYER_ID,
		PlayerWithTurnID: playerID, // since we've created an opponent, match can be started. TODO: remove
		WinnerID:         -1,
	}

	match.AddOrRefreshMatch(&matchState)
	updateMatchPlayerFields(&matchState)

	return &matchState, nil
}

func JoinMatch(playerID int, matchID common.Maybe[uuid.UUID], conn *websocket.Conn) {
	var matchState *models.Match
	var playerState *models.PlayerMatchState
	var opponentState *models.PlayerMatchState
	var err error

	if !matchID.HasValue {
		matchState, err = createNewMatchForPlayer(playerID, conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		playerState = matchState.Player0State.Value
		opponentState = matchState.Player1State.Value
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
			playerState = matchState.Player0State.Value
			opponentState = matchState.Player1State.Value

		} else if matchState.Player1State.HasValue && matchState.Player1State.Value.PlayerID == playerID {
			matchState.Player1State.Value.Connection = conn // updating connection just in case
			playerState = matchState.Player1State.Value
			opponentState = matchState.Player0State.Value

		} else if !matchState.Player0State.HasValue {
			// we are joining as first player
			newPlayerState, err := createInitialPlayerMatchState(playerID, conn)
			if err != nil {
				fmt.Println(err)
				return
			}
			matchState.Player0State = common.Maybe[models.PlayerMatchState]{
				HasValue: true,
				Value:    &newPlayerState,
			}

			if !matchState.Player1State.HasValue {
				matchState.PlayerWithTurnID = common.NONE_PLAYER_ID // match has only one player => noones turn
			} else {
				matchState.PlayerWithTurnID = selectRandomPlayer(matchState.Player0State.Value.PlayerID, matchState.Player1State.Value.PlayerID)
			}

			match.AddOrRefreshMatch(matchState)

			playerState = matchState.Player0State.Value
			opponentState = nil

		} else if !matchState.Player1State.HasValue {
			// we are joining as second player
			newPlayerState, err := createInitialPlayerMatchState(playerID, conn)
			if err != nil {
				fmt.Println(err)
				return
			}
			matchState.Player1State = common.Maybe[models.PlayerMatchState]{
				HasValue: true,
				Value:    &newPlayerState,
			}

			if !matchState.Player0State.HasValue {
				matchState.PlayerWithTurnID = common.NONE_PLAYER_ID // match has only one player => noones turn
			} else {
				matchState.PlayerWithTurnID = selectRandomPlayer(matchState.Player0State.Value.PlayerID, matchState.Player1State.Value.PlayerID)
			}

			match.AddOrRefreshMatch(matchState)

			playerState = matchState.Player1State.Value
			opponentState = matchState.Player0State.Value
		}

		updateMatchPlayerFields(matchState)
	}

	go senders.StartListeningBackendEvents(playerState, matchState)
	go senders.StartListeningBackendEvents(opponentState, matchState)

	senders.SendAllCardsToEveryone(matchState)
	senders.SendMatchInformationToEveryone(matchState)
	senders.SendAllCardInstancesToEveryone(matchState)
	playerState.Events <- enums.BackendEventHandChanged
	if opponentState != nil {
		opponentState.Events <- enums.BackendEventOpponentHandChanged
	}
	playerState.Events <- enums.BackendEventDeckChanged
	if opponentState != nil {
		opponentState.Events <- enums.BackendEventOpponentDeckChanged
	}
	playerState.Events <- enums.BackendEventDiscardPileChanged
	if opponentState != nil {
		opponentState.Events <- enums.BackendEventOpponentDiscardPileChanged
	}
}
