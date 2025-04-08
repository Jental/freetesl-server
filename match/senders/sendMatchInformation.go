package senders

import (
	"fmt"
	"log"

	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/db/queries"
	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

func SendMatchInformationToPlayer(playerState *models.PlayerMatchState, matchState *models.Match) error {
	if playerState.Connection == nil {
		return nil // Fake opponent has nil connection. TODO: the check should be removed
	}

	playerID := playerState.PlayerID
	opponentID := playerState.OpponentState.PlayerID
	playerIDs := []int{playerID, opponentID}
	players, err := queries.GetPlayersByIDs(playerIDs)
	if err != nil {
		return err
	}

	player, exists := players[playerID]
	if !exists {
		return fmt.Errorf("player with id '%d' is not found", playerID)
	}
	var opponent *dbModels.Player
	opponent, exists = players[opponentID]
	if !exists {
		return fmt.Errorf("player with id '%d' is not found", opponentID)
	}

	var dto = dtos.MatchInformationDTO{
		Player: &dtos.PlayerInformationDTO{
			ID:         playerID,
			Name:       player.DisplayName,
			AvatarName: player.AvatarName,
			State:      byte(enums.PlayerStateInMatch),
		},
		Opponent: &dtos.PlayerInformationDTO{
			ID:         opponentID,
			Name:       opponent.DisplayName,
			AvatarName: opponent.AvatarName,
			State:      byte(enums.PlayerStateInMatch),
		},
		HasRing:         playerState.HasRing(),
		OpponentHasRing: playerState.OpponentState.HasRing(),
	}

	var json = map[string]interface{}{
		"method": "matchInformationUpdate",
		"body":   dto,
	}

	log.Printf("[%d]: sending: matchInformationUpdate", playerState.PlayerID)

	// TODO: each active player should have two queues:
	// - of requests from client to be processed
	// - of messages from server
	//   ideally with some filtration to avoid sending multiple matchStates one after another
	err = sendJson(playerState, json)
	if err != nil {
		return err
	}

	log.Printf("[%d]: sent: matchInformationUpdate", playerState.PlayerID)

	return nil
}
