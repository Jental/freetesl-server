package interceptors

import (
	"fmt"

	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/models/enums"
)

type CoverInterceptor struct{}

func (ic CoverInterceptor) Execute(context *models.InterceptorContext) error {
	if context.TargetCardInstance == nil {
		return fmt.Errorf("[%d]: CoverInterceptor: no target card instance specified", context.PlayerState.PlayerID)
	}

	targetCreatureCardInstance, ok := context.TargetCardInstance.(*models.CardInstanceCreature)
	if !ok {
		return fmt.Errorf("[%d]: CoverInterceptor: target card instance is not a creature", context.PlayerState.PlayerID)
	}

	if targetCreatureCardInstance.HasEffect(enums.EffectTypeCover) {
		return fmt.Errorf("[%d]: CoverInterceptor: cover is present", context.PlayerState.PlayerID)
	}

	return nil
}
