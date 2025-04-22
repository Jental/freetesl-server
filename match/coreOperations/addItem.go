package coreOperations

import (
	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/models/enums"
)

func AddItem(playerState *models.PlayerMatchState, matchState *models.Match, cardInstance *models.CardInstanceCreature, itemCardInstance *models.CardInstanceItem) {
	cardInstance.Items = append(cardInstance.Items, itemCardInstance)
	itemCardInstance.EquippedTurnID = &matchState.TurnID

	playerState.SendEvent(enums.BackendEventCardInstancesChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentCardInstancesChanged)

	// to force lanes redraw
	playerState.SendEvent(enums.BackendEventLanesChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentLanesChanged)
}
