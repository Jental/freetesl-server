package actions

import (
	"fmt"
	"log"

	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/models/enums"
)

type ShackleAction struct{}

func (action ShackleAction) Execute(context *models.ActionContext) error {
	if context.CardID == nil {
		return fmt.Errorf("[%d]: ShackleAction: no CardID specified", context.PlayerState.PlayerID)
	}

	if context.TargetPlayerState == nil {
		return fmt.Errorf("[%d]: ShackleAction: no TargetPlayerState specified", context.PlayerState.PlayerID)
	}

	if context.TargetCardInstance == nil {
		return fmt.Errorf("[%d]: ShackleAction: no TargetCardInstance specified", context.PlayerState.PlayerID)
	}

	log.Printf("[%d]: ShackleAction; cardID: '%d'", context.PlayerState.PlayerID, *context.CardID)

	effect := models.Effect{
		EffectType:  enums.EffectTypeShackled,
		StartTurnID: context.PlayerState.MatchState.TurnID,
	}
	coreOperations.AddEffect(context.TargetPlayerState, context.TargetCardInstance, &effect)

	return nil
}
