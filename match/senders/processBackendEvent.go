package senders

import (
	"errors"
	"log"

	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

func ProcessBackendEvent(playerState *models.PlayerMatchState, event enums.BackendEventType) error {
	var err error = nil

	log.Printf("[%d]: event: '%s'(%d)\n", playerState.PlayerID, enums.BackendEventTypeName[event], event)

	// TODO: calculate dto hashes and don't do resent if dto has not changed
	//       + maybe some throttling will be a good idea
	//         collect events during some small interval, and send onlu unique messages after

	switch event {
	case enums.BackendEventDeckChanged, enums.BackendEventOpponentDeckChanged:
		err = SendDeckStateToPlayer(playerState, playerState.MatchState)
	case enums.BackendEventHealthChanged, enums.BackendEventManaChanged, enums.BackendEventHandChanged, enums.BackendEventLanesChanged, enums.BackendEventMatchStateRefresh,
		enums.BackendEventOpponentHealthChanged, enums.BackendEventOpponentManaChanged, enums.BackendEventOpponentHandChanged, enums.BackendEventOpponentLanesChanged, enums.BackendEventOpponentMatchStateRefresh,
		enums.BackendEventCardWatingForActionChanged, enums.BackendEventOpponentCardWatingForActionChanged,
		enums.BackendEventSwitchTurn:
		err = SendMatchStateToPlayer(playerState, playerState.MatchState)
	case enums.BackendEventDiscardPileChanged, enums.BackendEventOpponentDiscardPileChanged:
		err = SendDiscardPileStateToPlayer(playerState, playerState.MatchState)
	case enums.BackendEventCardInstancesChanged, enums.BackendEventOpponentCardInstancesChanged:
		err = SendAllCardInstancesToPlayer(playerState, playerState.MatchState)
	case enums.BackendEventMatchStart:
		err0 := SendMatchInformationToPlayer(playerState, playerState.MatchState)
		err1 := SendAllCardsToPlayer(playerState)
		if err0 != nil || err1 != nil {
			err = errors.Join(err0, err1)
		}
	case enums.BackendEventMatchEnd:
		err = SendMatchEndToPlayer(playerState, playerState.MatchState)
	}
	return err
}
