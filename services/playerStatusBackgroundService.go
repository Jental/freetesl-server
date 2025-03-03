package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/models/enums"
)

func StartPlayersActivityMonitoring() {
	for {
		var matchesToFinish map[uuid.UUID]*struct {
			hasWinner bool
			winnerID  int
		} = make(map[uuid.UUID]*struct {
			hasWinner bool
			winnerID  int
		})

		for playerID, playerInfo := range playersRunimeInfo {
			var diff = time.Since(*playerInfo.LastActivityTime)
			if diff.Seconds() > 30 {
				if playerInfo.State == enums.PlayerStateInMatch {
					var match, _, opponentState, err = match.GetCurrentMatchState(playerID)
					if err != nil {
						continue // no match
					}

					foundMatchToFinish, exists := matchesToFinish[match.Id]
					if exists {
						foundMatchToFinish.hasWinner = false // we already going to finish this match because of timeout of other player. both have timeouts => draw
						foundMatchToFinish.winnerID = -1
					} else {
						matchesToFinish[match.Id] = &struct {
							hasWinner bool
							winnerID  int
						}{
							hasWinner: true,
							winnerID:  opponentState.PlayerID,
						}
					}
				}

				playerInfo.State = enums.PlayerStateOffline
			} else if playerInfo.State == enums.PlayerStateInMatch {
				var _, _, _, err = match.GetCurrentMatchState(playerID)
				if err != nil {
					playerInfo.State = enums.PlayerStateOnline
				}
			}
		}

		for matchID, endInfo := range matchesToFinish {
			match.EndMatchByID(matchID, endInfo.winnerID)
		}

		time.Sleep(1 * time.Second)
	}
}
