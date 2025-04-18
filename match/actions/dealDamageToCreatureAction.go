package actions

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/match/models"
)

type DealDamageToCreatureAction struct{}

func (action DealDamageToCreatureAction) Execute(context *models.ActionContext) error {
	if context.CardID == nil {
		return fmt.Errorf("[%d]: DealDamageToCreatureAction: no CardID specified", context.PlayerState.PlayerID)
	}

	if context.TargetLane == nil {
		return fmt.Errorf("[%d]: DealDamageToCreatureAction: no TargetLaneID specified", context.PlayerState.PlayerID)
	}

	// parameter string is expected to be a single number
	if context.ParametersValues == nil {
		return fmt.Errorf("[%d]: DealDamageToCreatureAction: ParametersValues is expected to be set", context.PlayerState.PlayerID)
	}
	damage, err := strconv.Atoi(*context.ParametersValues)
	if err != nil {
		return fmt.Errorf("[%d]: DealDamageToCreatureAction: ParametersValues is expected to be a single number string", context.PlayerState.PlayerID)
	}

	targetCreatureCardInstance, ok := context.TargetCardInstance.(*models.CardInstanceCreature)
	if !ok {
		return fmt.Errorf("[%d]: DealDamageToCreatureAction: TargetCardInstance is expected to be a creature", context.PlayerState.PlayerID)
	}

	log.Printf("[%d]: DealDamageToCreatureAction; cardID: '%d'; parameters: '%s'", context.PlayerState.PlayerID, *context.CardID, *context.ParametersValues)

	coreOperations.ReduceCardHealth(context.TargetPlayerState, targetCreatureCardInstance, context.TargetLane, damage)

	return nil
}
