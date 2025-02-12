package actions

import (
	"time"

	"math/rand"

	"github.com/jental/freetesl-server/models"
)

var rnd rand.Rand = *rand.New(rand.NewSource(time.Now().UnixNano()))

func PlayRandomCards(playerState *models.PlayerMatchState2) {
	for {
		var cardWasPlayed bool = false

		for i, card := range playerState.Hand {
			if card.Cost <= playerState.Mana {
				var laneID = byte(rnd.Intn(2))
				err := MoveCardToLane(playerState, card, i, laneID) // TODO: check card type
				if err == nil {
					cardWasPlayed = true
					break
				}
			}
		}

		if !cardWasPlayed {
			break
		}
	}
}
