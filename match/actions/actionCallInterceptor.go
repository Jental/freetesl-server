package actions

import (
	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/models/enums"
)

type ActionCallInterceptor struct {
	ActionID               enums.ActionID
	CardID                 int
	ActionParametersValues *string
}

func (ic ActionCallInterceptor) Execute(context *models.InterceptorContext) error {
	if context.CardID != nil && *context.CardID != ic.CardID {
		return nil // for some other card
	}

	actionContext := models.NewActionContext(
		context.PlayerState,
		context.OpponentState,
		&ic.CardID,
		ic.ActionParametersValues,
		context.TargetPlayerState,
		context.TargetCardInstance,
		context.TargetLane,
	)
	return ExecuteAction(ic.ActionID, &actionContext)
}
