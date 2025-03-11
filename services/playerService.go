package services

import (
	"errors"
	"time"

	"github.com/jental/freetesl-server/db"
	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
	"github.com/samber/lo"
)

var playersRunimeInfo map[int]*models.PlayerRuntimeInformation = make(map[int]*models.PlayerRuntimeInformation)

func GetPlayers(onlyInGamePlayers bool) ([]*models.Player, error) {
	playersFromDB, err := db.GetPlayers()
	if err != nil {
		return nil, err
	}

	var players = lo.FilterMap(playersFromDB, func(p *dbModels.Player, _ int) (*models.Player, bool) {
		playerInfo, exists := playersRunimeInfo[p.ID]

		var state enums.PlayerState
		if !exists {
			state = enums.PlayerStateOffline
		} else {
			state = playerInfo.State
		}

		if onlyInGamePlayers && state == enums.PlayerStateOffline {
			return nil, false
		}

		player := models.Player{
			ID:          p.ID,
			DisplayName: p.DisplayName,
			AvatarName:  p.AvatarName,
			State:       state,
		}
		return &player, true
	})

	return players, nil
}

func GetPlayer(playerID int) (*models.Player, error) {
	playersFromDB, err := db.GetPlayersByIDs([]int{playerID})
	if err != nil {
		return nil, err
	}
	playerFromDB, exists := playersFromDB[playerID]
	if !exists {
		return nil, errors.New("player not found")
	}

	playerInfo, exists := playersRunimeInfo[playerFromDB.ID]

	var state enums.PlayerState
	if !exists {
		state = enums.PlayerStateOffline
	} else {
		state = playerInfo.State
	}

	player := models.Player{
		ID:          playerFromDB.ID,
		DisplayName: playerFromDB.DisplayName,
		AvatarName:  playerFromDB.AvatarName,
		State:       state,
	}
	return &player, nil
}

func SetPlayerState(playerID int, state enums.PlayerState) {
	if state == enums.PlayerStateOffline {
		delete(playersRunimeInfo, playerID)
	} else {
		var now = time.Now()
		playersRunimeInfo[playerID] = &models.PlayerRuntimeInformation{
			State:            state,
			LastActivityTime: &now,
		}
	}
}

func UpdatePlayerLastActivityTime(playerID int) {
	info, exists := playersRunimeInfo[playerID]
	if exists {
		var now = time.Now()
		info.LastActivityTime = &now
	} else {
		SetPlayerState(playerID, enums.PlayerStateOnline)
	}
}
