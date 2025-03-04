package senders

import (
	"fmt"
	"log"

	"github.com/jental/freetesl-server/db"
	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

func SendMatchInformationToPlayer(playerState *models.PlayerMatchState, matchState *models.Match) error {
	if playerState.Connection == nil {
		return nil // Fake opponent has nil connection. TODO: the check should be removed
	}

	var playerID = playerState.PlayerID
	opponentID, opponentExists, err := match.GetOpponentID(matchState, playerID)
	if err != nil {
		return err
	}

	var playerIDs []int
	if opponentExists {
		playerIDs = []int{playerID, opponentID}
	} else {
		playerIDs = []int{playerID}
	}
	players, err := db.GetPlayersByIDs(playerIDs)
	if err != nil {
		return err
	}

	player, exists := players[playerID]
	if !exists {
		return fmt.Errorf("player with id '%d' is not found", playerID)
	}
	var opponent *dbModels.Player
	if opponentExists {
		opponent, exists = players[opponentID]
		if !exists {
			return fmt.Errorf("player with id '%d' is not found", opponentID)
		}
	}

	var dto = dtos.MatchInformationDTO{
		Player: &dtos.PlayerInformationDTO{
			ID:         playerID,
			Name:       player.DisplayName,
			AvatarName: player.AvatarName,
			State:      byte(enums.PlayerStateInMatch),
		},
	}
	if opponentExists {
		dto.Opponent = &dtos.PlayerInformationDTO{
			ID:         opponentID,
			Name:       opponent.DisplayName,
			AvatarName: opponent.AvatarName,
			State:      byte(enums.PlayerStateInMatch),
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

	log.Printf("[%d]: sent: matchInformationUpdate", playerState.PlayerID)

	return nil
}
