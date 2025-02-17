package senders

import (
	"fmt"

	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

func StartListeningBackendEvents(playerState *models.PlayerMatchState, matchState *models.Match) {
	for {
		var event = <-playerState.Events

		fmt.Printf("event: [%d]: %s\n", playerState.PlayerID, enums.BackendEventTypeName[event])

		switch event {
		case enums.BackendEventDeckChanged, enums.BackendEventOpponentDeckChanged:
			SendDeckStateToPlayer(playerState, matchState)
		case enums.BackendEventHealthChanged, enums.BackendEventManaChanged, enums.BackendEventHandChanged, enums.BackendEventLanesChanged, enums.BackendEventMatchStateRefresh,
			enums.BackendEventOpponentHealthChanged, enums.BackendEventOpponentManaChanged, enums.BackendEventOpponentHandChanged, enums.BackendEventOpponentLanesChanged, enums.BackendEventOpponentMatchStateRefresh,
			enums.BackendEventSwitchTurn:
			SendMatchStateToPlayer(playerState, matchState)
		case enums.BackendEventDiscardPileChanged, enums.BackendEventOpponentDiscardPileChanged:
			SendDiscardPileStateToPlayer(playerState, matchState)
		case enums.BackendEventMatchEnd:
			SendMatchEndToPlayer(playerState, matchState)
		}
	}
}
