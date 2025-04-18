package interceptors

import (
	"fmt"

	dbEnums "github.com/jental/freetesl-server/db/enums"
	"github.com/jental/freetesl-server/match/models"
)

type GuardInterceptor struct{}

func (ic GuardInterceptor) Execute(context *models.InterceptorContext) error {
	if context.PlayerState == nil {
		return fmt.Errorf("[%d]: GuardInterceptor: no PlayerState specified", context.PlayerState.PlayerID)
	}
	if context.TargetPlayerState == nil {
		return fmt.Errorf("[%d]: GuardInterceptor: no TargetPlayerState specified", context.PlayerState.PlayerID)
	}
	if context.SourceLane == nil {
		return fmt.Errorf("[%d]: GuardInterceptor: no SourceLane specified", context.PlayerState.PlayerID)
	}

	if context.TargetCardInstance != nil { // hitCard
		if context.TargetCardInstance.HasKeyword(dbEnums.CardKeywordGuard) {
			return nil
		}
	}

	opponentLaneCards := context.OpponentState.GetLaneCards(context.SourceLane.Position)

	opponentGuardPresent := false
	for _, ocard := range opponentLaneCards {
		if ocard.HasKeyword(dbEnums.CardKeywordGuard) {
			opponentGuardPresent = true
			break
		}
	}

	if opponentGuardPresent {
		return fmt.Errorf("[%d]: GuardInterceptor: guard is present", context.PlayerState.PlayerID)
	}

	return nil
}
