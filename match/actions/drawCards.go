package actions

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/match/models"
)

type DrawCardsAction struct{}

func (action DrawCardsAction) Execute(context *models.ActionContext) error {
	if context.CardID == nil {
		return fmt.Errorf("[%d]: DrawCardsAction: no CardID specified", context.PlayerState.PlayerID)
	}

	// parameter string is expected to be a single number
	if context.ParametersValues == nil {
		return fmt.Errorf("[%d]: DrawCardsAction: ParametersValues is expected to be set", context.PlayerState.PlayerID)
	}
	amount, err := strconv.Atoi(*context.ParametersValues)
	if err != nil {
		return fmt.Errorf("[%d]: DrawCardsAction: ParametersValues is expected to be a single number string", context.PlayerState.PlayerID)
	}

	log.Printf("[%d]: DrawCardsActions; cardID: '%d'; parameters: '%s'", context.PlayerState.PlayerID, *context.CardID, *context.ParametersValues)

	for i := 0; i < amount; i++ {
		coreOperations.DrawCard(context.PlayerState)
	}

	return nil
}
