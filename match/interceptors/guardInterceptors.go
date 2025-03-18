package interceptors

import (
	"fmt"

	dbEnums "github.com/jental/freetesl-server/db/enums"
	"github.com/jental/freetesl-server/models"
)

type GuardInterceptor struct{}

func (ic GuardInterceptor) Execute(context *models.InterceptorContext) error {
	if context.SourceLaneID == nil {
		return fmt.Errorf("[%d]: GuardInterceptor: no lane id specified", context.PlayerState.PlayerID)
	}

	opponentLaneCards := context.OpponentState.GetLaneCards(*context.SourceLaneID)

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
