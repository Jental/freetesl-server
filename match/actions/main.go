package actions

import (
	"fmt"
	"log"

	"github.com/jental/freetesl-server/db/queries"
	"github.com/jental/freetesl-server/match/interceptors"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

var actions map[enums.ActionID]*models.Action = make(map[enums.ActionID]*models.Action)

func ExecuteAction(actionID enums.ActionID, context *models.ActionContext) error {
	action, exists := actions[actionID]
	if !exists {
		return fmt.Errorf("ExecuteAction: action with id '%s' does not exist/registered", actionID)
	}

	return (*action).Execute(context)
}

// all existing actions are expected to be registered => registering them here, not in the main file
func RegisterAllActions() {
	var dealDamageToCreatureAction models.Action = DealDamageToCreatureAction{}
	actions[enums.ActionIDDealDamageToCreature] = &dealDamageToCreatureAction

	var drawCardsAction models.Action = DrawCardsAction{}
	actions[enums.ActionIDDrawCards] = &drawCardsAction

	var shackleAction models.Action = ShackleAction{}
	actions[enums.ActionShackle] = &shackleAction

	var healAction models.Action = HealAction{}
	actions[enums.ActionHeal] = &healAction
}

func RegisterActionsForCards() {
	actionsFromDB, err := queries.GetCardActions()
	if err != nil {
		log.Panic(err)
		return
	}

	for _, action := range actionsFromDB {
		actionID := enums.ActionID(action.ActionID)
		_, exists := actions[actionID]
		if !exists {
			log.Println(fmt.Errorf("RegisterActionsFromDB: action with id '%s' does not exist/registered", actionID))
			continue // TODO: panic and return
		}

		interceptorPointID := enums.InteceptorPoint(action.InterceptorPointID)
		var interceptor models.Interceptor = ActionCallInterceptor{ActionID: actionID, CardID: action.CardID, ActionParametersValues: action.ActionParametersValues}
		interceptors.RegisterInterceptor(interceptorPointID, &interceptor)
	}
}
