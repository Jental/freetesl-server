package handlers

import (
	"github.com/jental/freetesl-server/match/senders"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

func JoinMatch(playerState *models.PlayerMatchState) {
	senders.SendAllCardsToPlayer(playerState)

	playerState.Events <- enums.BackendEventMatchStart
	playerState.Events <- enums.BackendEventCardInstancesChanged
	playerState.Events <- enums.BackendEventOpponentCardInstancesChanged
	playerState.Events <- enums.BackendEventHandChanged
	playerState.Events <- enums.BackendEventOpponentHandChanged
	playerState.Events <- enums.BackendEventDeckChanged
	playerState.Events <- enums.BackendEventOpponentDeckChanged
	playerState.Events <- enums.BackendEventDiscardPileChanged
	playerState.Events <- enums.BackendEventOpponentDiscardPileChanged
}
