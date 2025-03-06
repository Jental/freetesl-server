package services

import (
	"log"
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

		var playersToDeleteRuntimeInfo = make([]int, 0)

		for playerID, playerInfo := range playersRunimeInfo {
			var diff = time.Since(*playerInfo.LastActivityTime)
			if diff.Seconds() > 120 {
				log.Printf("[%d]: got inactive", playerID)
				if playerInfo.State == enums.PlayerStateInMatch {
					// TODO: better match timeouts handling (based on turn timeouts)
					var match, _, opponentState, err = match.GetCurrentMatchState(playerID)
					if err != nil {
						continue // no match
					}

					log.Printf("[%d]: got inactive: match: %s", playerID, match.Id.String())

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
				playersToDeleteRuntimeInfo = append(playersToDeleteRuntimeInfo, playerID)
			} else if playerInfo.State == enums.PlayerStateInMatch {
				var _, _, _, err = match.GetCurrentMatchState(playerID)
				if err != nil {
					log.Printf("[%d]: match seems to be finished - changing state to 'online'", playerID)
					playerInfo.State = enums.PlayerStateOnline
				}
			}
		}

		for matchID, endInfo := range matchesToFinish {
			match.EndMatchByID(matchID, endInfo.winnerID)
		}

		for _, playerID := range playersToDeleteRuntimeInfo {
			log.Printf("[%d]: removing player runtime info", playerID)
			delete(playersRunimeInfo, playerID)
		}

		time.Sleep(1 * time.Second)
	}
}
