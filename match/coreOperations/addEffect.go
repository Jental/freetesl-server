package coreOperations

import (
	"fmt"

	"github.com/jental/freetesl-server/match/effects"
	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/models/enums"
)

func AddEffect(playerState *models.PlayerMatchState, cardInstance models.CardInstance, effect *effects.EffectInstance) error {
	creatureCardInstance, castSuccessed := cardInstance.(*models.CardInstanceCreature)
	if castSuccessed {
		creatureCardInstance.Effects = append(creatureCardInstance.Effects, effect)
	} else {
		return fmt.Errorf("AddEffect: effects can only be added to creatures and items")
	}

	playerState.SendEvent(enums.BackendEventCardInstancesChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentCardInstancesChanged)

	// to force lanes redraw
	playerState.SendEvent(enums.BackendEventLanesChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentLanesChanged)

	return nil
}
