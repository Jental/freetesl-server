package interceptors

import (
	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

// TODO: I think that's not the only case, when we should do something, when a card is present on one of lanes
// Something like actionsInterceptor, but one calling action, when card is present
type BrumaProfiteerInterceptor struct{}

func (ic BrumaProfiteerInterceptor) Execute(context *models.InterceptorContext) error {
	laneCards := context.PlayerState.GetAllLaneCards()

	cardIsPresent := false
	for _, ocard := range laneCards {
		if ocard.Card.ID == enums.CardBrumaProfiteer && ocard.CardInstanceID != context.CardInstanceID {
			cardIsPresent = true
			break
		}
	}

	if cardIsPresent {
		coreOperations.IncreasePlayerHealth(context.PlayerState, 1)
	}

	return nil
}
