package main

import (
	"errors"
	"fmt"
	"maps"
	"math/rand"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/mappers"
	"github.com/jental/freetesl-server/models"
	"github.com/samber/lo"
)

var matches map[uuid.UUID]*models.Match = nil
var playersToMatches map[int]*models.Match = nil
var once sync.Once
var mutex sync.Mutex
var rnd rand.Rand = *rand.New(rand.NewSource(time.Now().UnixNano()))

func createMatchesIfNeeded() {
	once.Do(func() {
		var m = make(map[uuid.UUID]*models.Match)
		matches = m

		var p = make(map[int]*models.Match)
		playersToMatches = p
	})
}

func getMatch(matchID uuid.UUID) (*models.Match, bool) {
	createMatchesIfNeeded()

	mutex.Lock()
	defer mutex.Unlock()

	match, exist := matches[matchID]
	return match, exist
}

func addMatch(match *models.Match, playerID int) {
	createMatchesIfNeeded()

	mutex.Lock()
	defer mutex.Unlock()

	matches[match.Id] = match
	playersToMatches[playerID] = match
}

func getCurrentMatchState(playerID int) (*models.Match, *models.PlayerMatchState2, error) {
	match, exist := playersToMatches[playerID]
	if !exist {
		return nil, nil, fmt.Errorf("player with id '%d' have no active match", playerID)
	}

	if !match.Player0State.HasValue || !match.Player1State.HasValue {
		return nil, nil, fmt.Errorf("match for player '%d' is not started yet - waiting for second party", playerID)
	}

	var playerState *models.PlayerMatchState2
	if match.Player0State.Value.PlayerID == playerID {
		playerState = match.Player0State.Value
	} else if match.Player1State.Value.PlayerID == playerID {
		playerState = match.Player1State.Value
	} else {
		return nil, nil, fmt.Errorf("match for player '%d' is started for different players. this should not happen", playerID)
	}

	return match, playerState, nil
}

func createInitialPlayerMatchState(playerID int, conn *websocket.Conn) (models.PlayerMatchState2, error) {
	decks, err := getDecks(playerID)
	if err != nil {
		return models.PlayerMatchState2{}, err
	}

	var deckInstance []*models.CardInstance = lo.Shuffle(
		lo.FlatMap(
			slices.Collect(maps.Values(decks[0].Cards)),
			func(cardWithCount models.CardWithCount, _ int) []*models.CardInstance {
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

func selectRandomPlayer(player0ID int, player1ID int) int {
	var idx = rnd.Intn(2)
	if idx == 0 {
		return player0ID
	} else {
		return player1ID
	}
}

func joinMatch(playerID int, matchID common.Maybe[uuid.UUID], conn *websocket.Conn) {
	var match models.Match
	if !matchID.HasValue {
		playerState, err := createInitialPlayerMatchState(playerID, conn)
		if err != nil {
			fmt.Println(err)
			return
		}

		// TODO: for debug, shoud be removed later
		// this block auocreates opponent
		var opponentID = 2
		player1State, err := createInitialPlayerMatchState(opponentID, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		match = models.Match{
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
		}

		addMatch(&match, playerID)
	} else {
		match, exist := getMatch(*matchID.Value)
		if !exist {
			fmt.Println(errors.New("match with given id does not exist"))
			return
		} else if match.Player0State.HasValue && match.Player1State.HasValue && match.Player0State.Value.PlayerID != playerID && match.Player1State.Value.PlayerID != playerID {
			fmt.Println(errors.New("match is already started with different players"))
		} else if match.Player0State.HasValue && match.Player0State.Value.PlayerID == playerID {
			match.Player0State.Value.Connection = conn // updating connection just in case
			playersToMatches[playerID] = match         // updating just in case too
		} else if match.Player1State.HasValue && match.Player1State.Value.PlayerID == playerID {
			match.Player1State.Value.Connection = conn // updating connection just in case
			playersToMatches[playerID] = match         // updating just in case too
		} else if !match.Player0State.HasValue {
			// we are joining as first player
			playerState, err := createInitialPlayerMatchState(playerID, conn)
			if err != nil {
				fmt.Println(err)
				return
			}
			match.Player0State = common.Maybe[models.PlayerMatchState2]{
				HasValue: true,
				Value:    &playerState,
			}

			playersToMatches[playerID] = match

			if !match.Player1State.HasValue {
				match.PlayerWithTurnID = common.NONE_PLAYER_ID // match has only one player => noones turn
			} else {
				match.PlayerWithTurnID = selectRandomPlayer(match.Player0State.Value.PlayerID, match.Player1State.Value.PlayerID)
			}
		} else if !match.Player1State.HasValue {
			// we are joining as second player
			playerState, err := createInitialPlayerMatchState(playerID, conn)
			if err != nil {
				fmt.Println(err)
				return
			}
			match.Player1State = common.Maybe[models.PlayerMatchState2]{
				HasValue: true,
				Value:    &playerState,
			}

			playersToMatches[playerID] = match

			if !match.Player0State.HasValue {
				match.PlayerWithTurnID = common.NONE_PLAYER_ID // match has only one player => noones turn
			} else {
				match.PlayerWithTurnID = selectRandomPlayer(match.Player0State.Value.PlayerID, match.Player1State.Value.PlayerID)
			}
		}
	}

	sendMatchInformationToEveryone(&match)
	sendMatchStateToEveryone(&match)
}

func getOpponentID(match *models.Match, playerID int) (int, bool, error) {
	if match.Player0State.HasValue && match.Player1State.HasValue {
		if match.Player0State.Value.PlayerID == playerID {
			return match.Player1State.Value.PlayerID, true, nil
		} else if match.Player1State.Value.PlayerID == playerID {
			return match.Player0State.Value.PlayerID, true, nil
		} else {
			return -1, false, fmt.Errorf("player with id '%d' is not a part of a match", playerID)
		}
	} else if match.Player0State.HasValue {
		if match.Player0State.Value.PlayerID == playerID {
			return -1, false, nil
		} else {
			return -1, false, fmt.Errorf("player with id '%d' is not a part of a match", playerID)
		}
	} else if match.Player1State.HasValue {
		if match.Player1State.Value.PlayerID == playerID {
			return -1, false, nil
		} else {
			return -1, false, fmt.Errorf("player with id '%d' is not a part of a match", playerID)
		}
	} else {
		return -1, false, fmt.Errorf("player with id '%d' is not a part of a match", playerID)
	}
}

func sendMatchInformationToEveryone(match *models.Match) {
	if match.Player0State.HasValue {
		go sendMatchInformationToPlayerWithErrorHandling(match.Player0State.Value, match)
	}
	if match.Player0State.HasValue {
		go sendMatchInformationToPlayerWithErrorHandling(match.Player1State.Value, match)
	}
}

func sendMatchInformationToPlayerWithErrorHandling(playerState *models.PlayerMatchState2, match *models.Match) {
	var err = sendMatchInformationToPlayer(playerState, match)
	if err != nil {
		fmt.Println(err)
	}
}

func sendMatchInformationToPlayer(playerState *models.PlayerMatchState2, match *models.Match) error {
	if playerState.Connection == nil {
		return nil // Fake opponent has nil connection. TODO: the check should be removed
	}

	var playerID = playerState.PlayerID
	opponentID, opponentExists, err := getOpponentID(match, playerID)
	if err != nil {
		return err
	}

	var playerIDs []int
	if opponentExists {
		playerIDs = []int{playerID, opponentID}
	} else {
		playerIDs = []int{playerID}
	}
	players, err := getPlayers(playerIDs)
	if err != nil {
		return err
	}

	player, exists := players[playerID]
	if !exists {
		return fmt.Errorf("player with id '%d' is not found", playerID)
	}
	var opponent *models.Player
	if opponentExists {
		opponent, exists = players[opponentID]
		if !exists {
			return fmt.Errorf("player with id '%d' is not found", opponentID)
		}
	}

	var dto = dtos.MatchInformationDTO{
		Player: &dtos.PlayerInformationDTO{
			Name:       player.DisplayName,
			AvatarName: player.AvatarName,
		},
	}
	if opponentExists {
		dto.Opponent = &dtos.PlayerInformationDTO{
			Name:       opponent.DisplayName,
			AvatarName: opponent.AvatarName,
		}
	} else {
		dto.Opponent = nil
	}

	var json = map[string]interface{}{
		"method": "matchInformationUpdate",
		"body":   dto,
	}

	// TODO: each active player should have two queues:
	// - of requests from client to be processed
	// - of messages from server
	//   ideally with some filtration to avoid sending multiple matchStates one after another
	err = playerState.Connection.WriteJSON(json)
	if err != nil {
		return err
	}

	return nil
}

func sendMatchStateToEveryone(match *models.Match) {
	if match.Player0State.HasValue {
		go sendMatchStateToPlayerWithErrorHandling(match.Player0State.Value, match)
	}

	if match.Player1State.HasValue {
		go sendMatchStateToPlayerWithErrorHandling(match.Player1State.Value, match)
	}
}

func sendMatchStateToPlayerWithErrorHandling(playerState *models.PlayerMatchState2, match *models.Match) {
	var err = sendMatchStateToPlayer(playerState, match)
	if err != nil {
		fmt.Println(err)
	}
}

func sendMatchStateToPlayer(playerState *models.PlayerMatchState2, match *models.Match) error {
	if playerState.Connection == nil {
		return nil // Fake opponent has nil connection. TODO: the check should be removed
	}

	var dto, err = mappers.MapToMatchStateDTO(match, playerState.PlayerID)
	if err != nil {
		return err
	}
	var json = map[string]interface{}{
		"method": "matchStateUpdate",
		"body":   dto,
	}

	// TODO: each active player should have two queues:
	// - of requests from client to be processed
	// - of messages from server
	//   ideally with some filtration to avoid sending multiple matchStates one after another
	err = playerState.Connection.WriteJSON(json)
	if err != nil {
		return err
	}

	return nil
}

func drawCard(playerState *models.PlayerMatchState2) {
	if len(playerState.Deck) == 0 {
		return // TODO: for now doing nothing, but later next rune should be broken
	}

	var drawnCard = playerState.Deck[0]
	drawnCard.IsActive = true
	playerState.Hand = append(playerState.Hand, drawnCard)
	playerState.Deck = playerState.Deck[1:]
}

func switchTurn(match *models.Match) {
	var isFirstPlayersTurn = match.Player0State.Value.PlayerID == match.PlayerWithTurnID
	if isFirstPlayersTurn {
		match.PlayerWithTurnID = match.Player1State.Value.PlayerID
	} else {
		match.PlayerWithTurnID = match.Player0State.Value.PlayerID
	}
}

func endTurn(playerID int) {
	match, playerState, err := getCurrentMatchState(playerID)
	if err != nil {
		fmt.Println(err)
		return
	}

	switchTurn(match)
	sendMatchStateToEveryone(match)

	time.Sleep(3 * time.Second)
	switchTurn(match)
	drawCard(playerState)

	for _, card := range playerState.LeftLaneCards {
		card.IsActive = true
	}

	playerState.MaxMana = playerState.MaxMana + 1
	playerState.Mana = playerState.MaxMana

	sendMatchStateToEveryone(match)
}

func moveCardToLane(playerID int, cardInstanceID uuid.UUID, laneID byte) {
	match, playerState, err := getCurrentMatchState(playerID)
	if err != nil {
		fmt.Println(err)
		return
	}

	var idx = slices.IndexFunc(playerState.Hand, func(el *models.CardInstance) bool { return el.CardInstanceID == cardInstanceID })
	if idx < 0 {
		fmt.Println(fmt.Errorf("card instance with id '%s' is not present in a hand of a player '%d'", cardInstanceID, playerID))
		return
	}
	var cardInstance = playerState.Hand[idx]

	if cardInstance.Cost > playerState.Mana {
		fmt.Println(fmt.Errorf("not enough mana '%d' of '%d'", cardInstance.Cost, playerState.Mana))
		sendMatchStateToEveryone(match)
		return
	}

	if laneID == common.LEFT_LANE_ID {
		if len(playerState.LeftLaneCards) >= common.MAX_LANE_CARDS {
			fmt.Println(fmt.Errorf("lane is already full"))
			return
		}
		playerState.LeftLaneCards = append(playerState.LeftLaneCards, cardInstance)
	} else if laneID == common.RIGHT_LANE_ID {
		// if len(playerState.RightLaneCards) >= common.MAX_LANE_CARDS {
		// 	fmt.Println(fmt.Errorf("lane is already full"))
		// 	return
		// }
		// playerState.RightLaneCards = append(playerState.LeftLaneCards, cardInstance)
	} else {
		fmt.Println(fmt.Errorf("invali lane id: %d", laneID))
		return
	}

	cardInstance.IsActive = false
	playerState.Hand = slices.Delete(playerState.Hand, idx, idx+1)
	playerState.Mana = playerState.Mana - cardInstance.Cost

	sendMatchStateToEveryone(match)
}
