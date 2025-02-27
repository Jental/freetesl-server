package services

import (
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/jental/freetesl-server/db"
	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
	"github.com/samber/lo"
)

var playersRunimeInfo map[int]*models.PlayerRuntimeInformation = make(map[int]*models.PlayerRuntimeInformation)

func GetPlayers() ([]*models.Player, error) {
	playersFromDB, err := db.GetPlayers(nil)
	if err != nil {
		return nil, err
	}

	var players = lo.Map[*dbModels.Player, *models.Player](slices.Collect(maps.Values(playersFromDB)), func(p *dbModels.Player, _ int) *models.Player {
		playerInfo, exists := playersRunimeInfo[p.ID]

		var state enums.PlayerState
		if !exists {
			state = enums.PlayerStateOffline
		} else {
			if playerInfo.LastActivityTime == nil {
				fmt.Println(fmt.Errorf("unexpected null LastActivityTime"))
				state = enums.PlayerStateOffline

			} else {
				var diff = time.Since(*playerInfo.LastActivityTime)
				if diff.Seconds() > 30 {
					state = enums.PlayerStateOffline
				} else {
					state = playerInfo.State
				}
			}
		}

		player := models.Player{
			ID:          p.ID,
			DisplayName: p.DisplayName,
			AvatarName:  p.AvatarName,
			State:       state,
		}
		return &player
	})

	return players, nil
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
