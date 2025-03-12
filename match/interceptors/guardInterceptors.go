package interceptors

import (
	"fmt"

	dbEnums "github.com/jental/freetesl-server/db/enums"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

type GuardInterceptor struct{}

func (ic GuardInterceptor) GetInterceptorPoint() enums.InteceptorPoint {
	return enums.InterceptorPointHitFaceBefore
}

func (ic GuardInterceptor) Execute(context *models.InterceptorContext) error {
	if context.LaneID == nil {
		return fmt.Errorf("[%d]: GuardInterceptor: no lane id specified", context.PlayerState.PlayerID)
	}

	opponentLaneCards := context.OpponentState.GetLaneCards(*context.LaneID)

	opponentGuardPresent := false
OuterLoop:
	for _, ocard := range opponentLaneCards {
		for _, kw := range ocard.Card.Keywords {
			if kw == dbEnums.CardKeywordGuard {
				opponentGuardPresent = true
				break OuterLoop
			}
		}
	}

	if opponentGuardPresent {
		return fmt.Errorf("[%d]: GuardInterceptor: guard is present", context.PlayerState.PlayerID)
	}

	return nil
}
