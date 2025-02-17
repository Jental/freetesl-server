package senders

import (
	"fmt"

	"github.com/jental/freetesl-server/db"
	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/models"
)

func SendMatchInformationToEveryone(match *models.Match) {
	if match.Player0State.HasValue {
		go sendMatchInformationToPlayerWithErrorHandling(match.Player0State.Value, match)
	}
	if match.Player0State.HasValue {
		go sendMatchInformationToPlayerWithErrorHandling(match.Player1State.Value, match)
	}
}

func sendMatchInformationToPlayerWithErrorHandling(playerState *models.PlayerMatchState, match *models.Match) {
	var err = sendMatchInformationToPlayer(playerState, match)
	if err != nil {
		fmt.Println(err)
	}
}

func sendMatchInformationToPlayer(playerState *models.PlayerMatchState, matchState *models.Match) error {
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
	players, err := db.GetPlayers(playerIDs)
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

	fmt.Printf("sent: [%d]: matchInformationUpdate\n", playerState.PlayerID)

	return nil
}
