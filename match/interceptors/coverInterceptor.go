package interceptors

import (
	"fmt"

	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

type CoverInterceptor struct{}

func (ic CoverInterceptor) Execute(context *models.InterceptorContext) error {
	if context.TargetCardInstance == nil {
		return fmt.Errorf("[%d]: CoverInterceptor: no target card instance specified", context.PlayerState.PlayerID)
	}

	if context.TargetCardInstance.HasEffect(enums.EffectTypeCover) {
		return fmt.Errorf("[%d]: CoverInterceptor: cover is present", context.PlayerState.PlayerID)
	}

	return nil
}
