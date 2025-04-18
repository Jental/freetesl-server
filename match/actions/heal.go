package actions

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/match/models"
)

type HealAction struct{}

func (action HealAction) Execute(context *models.ActionContext) error {
	if context.CardID == nil {
		return fmt.Errorf("[%d]: HealAction: no CardID specified", context.PlayerState.PlayerID)
	}

	amount, err := strconv.Atoi(*context.ParametersValues)
	if err != nil {
		return fmt.Errorf("[%d]: HealAction: ParametersValues is expected to be a single number string", context.PlayerState.PlayerID)
	}

	log.Printf("[%d]: HealAction; cardID: '%d'; parameters: '%s'", context.PlayerState.PlayerID, *context.CardID, *context.ParametersValues)

	coreOperations.IncreasePlayerHealth(context.PlayerState, amount)

	return nil
}
