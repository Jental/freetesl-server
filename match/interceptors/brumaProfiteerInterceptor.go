package interceptors

import (
	"fmt"

	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/models/enums"
)

// TODO: I think that's not the only case, when we should do something, when a card is present on one of lanes
// Something like actionsInterceptor, but one calling action, when card is present
type BrumaProfiteerInterceptor struct{}

func (ic BrumaProfiteerInterceptor) Execute(context *models.InterceptorContext) error {
	if context.CardInstanceID == nil {
		return fmt.Errorf("[%d]: BrumaProfiteerInterceptor: no target card instance id specified", context.PlayerState.PlayerID)
	}

	laneCards := context.PlayerState.GetAllLaneCardInstances()

	cardIsPresent := false
	for _, ocard := range laneCards {
		if ocard.Card.ID == enums.CardBrumaProfiteer && ocard.CardInstanceID != *context.CardInstanceID {
			cardIsPresent = true
			break
		}
	}

	if cardIsPresent {
		coreOperations.IncreasePlayerHealth(context.PlayerState, 1)
	}

	return nil
}
