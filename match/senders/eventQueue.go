package senders

import (
	"log"

	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

func StartListeningBackendEvents(playerState *models.PlayerMatchState, matchState *models.Match) {
	for {
		var event = <-playerState.Events
		var err error = nil

		log.Printf("[%d]: event: %s\n", playerState.PlayerID, enums.BackendEventTypeName[event])

		// TODO: calculate dto hashes and don't do resent if dto has not changed
		//       + maybe some throttling will be a good ided
		//         collect events during some small interval, and send onlu unique messages after

		var stopListening = false

		switch event {
		case enums.BackendEventDeckChanged, enums.BackendEventOpponentDeckChanged:
			err = SendDeckStateToPlayer(playerState, matchState)
		case enums.BackendEventHealthChanged, enums.BackendEventManaChanged, enums.BackendEventHandChanged, enums.BackendEventLanesChanged, enums.BackendEventMatchStateRefresh,
			enums.BackendEventOpponentHealthChanged, enums.BackendEventOpponentManaChanged, enums.BackendEventOpponentHandChanged, enums.BackendEventOpponentLanesChanged, enums.BackendEventOpponentMatchStateRefresh,
			enums.BackendEventSwitchTurn:
			err = SendMatchStateToPlayer(playerState, matchState)
		case enums.BackendEventDiscardPileChanged, enums.BackendEventOpponentDiscardPileChanged:
			err = SendDiscardPileStateToPlayer(playerState, matchState)
		case enums.BackendEventCardInstancesChanged, enums.BackendEventOpponentCardInstancesChanged:
			err = SendAllCardInstancesToPlayer(playerState, matchState)
		case enums.BackendEventMatchStart:
			err = SendMatchInformationToPlayer(playerState, matchState)
		case enums.BackendEventMatchEnd:
			err = SendMatchEndToPlayer(playerState, matchState)
			stopListening = true // we stop listening on match end
		}

		if err != nil {
			log.Printf("[%d]: sending error: '%s'", playerState.PlayerID, err)
		}

		if stopListening {
			break
		}
	}
}
