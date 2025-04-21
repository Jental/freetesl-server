package coreOperations

import (
	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/models/enums"
)

func AddItem(playerState *models.PlayerMatchState, cardInstance *models.CardInstanceCreature, itemCardInstance *models.CardInstanceItem) {
	cardInstance.Items = append(cardInstance.Items, itemCardInstance)

	playerState.SendEvent(enums.BackendEventCardInstancesChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentCardInstancesChanged)

	// to force lanes redraw
	playerState.SendEvent(enums.BackendEventLanesChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentLanesChanged)
}
