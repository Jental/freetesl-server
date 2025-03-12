package match

import (
	"fmt"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

var matches map[uuid.UUID]*models.Match = nil
var playersToMatches map[int]*models.Match = nil
var once sync.Once
var mutex sync.Mutex // TODO: store mutex in match object

func createMatchesIfNeeded() {
	once.Do(func() {
		var m = make(map[uuid.UUID]*models.Match)
		matches = m

		var p = make(map[int]*models.Match)
		playersToMatches = p
	})
}

func GetMatch(matchID uuid.UUID) (*models.Match, bool) {
	createMatchesIfNeeded()

	mutex.Lock()
	defer mutex.Unlock()

	match, exist := matches[matchID]
	return match, exist
}

func AddOrRefreshMatch(match *models.Match) {
	createMatchesIfNeeded()

	mutex.Lock()
	defer mutex.Unlock()

	matches[match.Id] = match

	if match.Player0State.HasValue {
		playersToMatches[match.Player0State.Value.PlayerID] = match
	}

	if match.Player1State.HasValue {
		playersToMatches[match.Player1State.Value.PlayerID] = match
	}
}

func EndMatch(match *models.Match, winnerID int) {
	log.Printf("EndMatch: %s", match.Id)

	match.WinnerID = winnerID

	match.Player0State.Value.SendEvent(enums.BackendEventMatchEnd)
	match.Player1State.Value.SendEvent(enums.BackendEventMatchEnd)

	go func() {
		time.Sleep(2 * time.Second) // to let some (at least matchEnd) events to be sent to FE
		DisconnectFromMatch(match.Player0State.Value)
		DisconnectFromMatch(match.Player1State.Value)
	}()

	mutex.Lock() // TODO: store mutex in match object
	defer mutex.Unlock()

	delete(playersToMatches, match.Player0State.Value.PlayerID)
	delete(playersToMatches, match.Player1State.Value.PlayerID)
	delete(matches, match.Id)
}

func EndMatchByID(matchID uuid.UUID, winnerID int) {
	match, exists := matches[matchID]
	if exists {
		EndMatch(match, winnerID)
	}
}

func GetCurrentMatchState(playerID int) (*models.Match, *models.PlayerMatchState, *models.PlayerMatchState, error) {
	match, exist := playersToMatches[playerID]
	if !exist {
		return nil, nil, nil, fmt.Errorf("player with id '%d' have no active match", playerID)
	}

	if !match.Player0State.HasValue || !match.Player1State.HasValue {
		return nil, nil, nil, fmt.Errorf("match for player '%d' is not started yet - waiting for second party", playerID)
	}

	var playerState, opponentState *models.PlayerMatchState
	if match.Player0State.Value.PlayerID == playerID {
		playerState = match.Player0State.Value
		opponentState = match.Player1State.Value
	} else if match.Player1State.Value.PlayerID == playerID {
		playerState = match.Player1State.Value
		opponentState = match.Player0State.Value
	} else {
		return nil, nil, nil, fmt.Errorf("match for player '%d' is started for different players. this should not happen", playerID)
	}

	return match, playerState, opponentState, nil
}

func GetOpponent(match *models.Match, playerID int) (*models.PlayerMatchState, bool, error) {
	if match.Player0State.HasValue && match.Player1State.HasValue {
		if match.Player0State.Value.PlayerID == playerID {
			return match.Player1State.Value, true, nil
		} else if match.Player1State.Value.PlayerID == playerID {
			return match.Player0State.Value, true, nil
		} else {
			return nil, false, fmt.Errorf("player with id '%d' is not a part of a match", playerID)
		}
	} else if match.Player0State.HasValue {
		if match.Player0State.Value.PlayerID == playerID {
			return nil, false, nil
		} else {
			return nil, false, fmt.Errorf("player with id '%d' is not a part of a match", playerID)
		}
	} else if match.Player1State.HasValue {
		if match.Player1State.Value.PlayerID == playerID {
			return nil, false, nil
		} else {
			return nil, false, fmt.Errorf("player with id '%d' is not a part of a match", playerID)
		}
	} else {
		return nil, false, fmt.Errorf("player with id '%d' is not a part of a match", playerID)
	}
}

func GetOpponentID(match *models.Match, playerID int) (int, bool, error) {
	var opponent, exist, err = GetOpponent(match, playerID)
	if err != nil {
		return -1, false, err
	}
	if !exist {
		return -1, false, nil
	}
	return opponent.PlayerID, true, nil
}

func GetCardInstanceFromHand(playerState *models.PlayerMatchState, cardInstanceID uuid.UUID) (*models.CardInstance, int, error) {
	var idx = slices.IndexFunc(playerState.GetHand(), func(el *models.CardInstance) bool { return el.CardInstanceID == cardInstanceID })
	if idx < 0 {
		return nil, -1, fmt.Errorf("card instance with id '%s' is not present in a hand of a player '%d'", cardInstanceID, playerState.PlayerID)
	}
	return playerState.GetHand()[idx], idx, nil
}

func GetCardInstanceFromLanes(playerState *models.PlayerMatchState, cardInstanceID uuid.UUID) (*models.CardInstance, enums.Lane, int, error) {
	var idx = slices.IndexFunc(playerState.GetLeftLaneCards(), func(el *models.CardInstance) bool { return el.CardInstanceID == cardInstanceID })
	if idx >= 0 {
		return playerState.GetLeftLaneCards()[idx], enums.LaneLeft, idx, nil
	}

	idx = slices.IndexFunc(playerState.GetRightLaneCards(), func(el *models.CardInstance) bool { return el.CardInstanceID == cardInstanceID })
	if idx >= 0 {
		return playerState.GetRightLaneCards()[idx], enums.LaneRight, idx, nil
	}

	return nil, 0, -1, fmt.Errorf("card instance with id '%s' is not present on lanes of a player '%d'", cardInstanceID, playerState.PlayerID)
}
