package operations

import (
	"fmt"

	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/match/interceptors"
	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/models/enums"
)

// TODO: implement with inteface after signature becomes clear
func hitCardCheck(playerState *models.PlayerMatchState, cardInstance *models.CardInstanceCreature, lane *models.Lane, opponentLane *models.Lane) error {
	if !cardInstance.IsActive() {
		return fmt.Errorf("[%d]: card with id '%s' is not active", playerState.PlayerID, cardInstance.CardInstanceID)
	}

	if lane.Position != opponentLane.Position {
		return fmt.Errorf("[%d]: cards are on different lanes", playerState.PlayerID)
	}

	return nil
}

// logic itself
func hitCard(
	playerState *models.PlayerMatchState,
	opponentState *models.PlayerMatchState,
	cardInstance *models.CardInstanceCreature,
	opponentCardInstance *models.CardInstanceCreature,
	lane *models.Lane,
	opponentLane *models.Lane,
) {
	coreOperations.ReduceCardHealth(opponentState, opponentCardInstance, opponentLane, cardInstance.GetComputedPower())
	coreOperations.ReduceCardHealth(playerState, cardInstance, lane, opponentCardInstance.GetComputedPower())
	cardInstance.SetIsActive(false)
}

func HitCard(
	playerState *models.PlayerMatchState,
	opponentState *models.PlayerMatchState,
	cardInstance *models.CardInstanceCreature,
	opponentCardInstance *models.CardInstanceCreature,
	lane *models.Lane,
	opponentLane *models.Lane,
) error {
	err := hitCardCheck(playerState, cardInstance, lane, opponentLane)
	if err != nil {
		return err
	}

	var opponentCardInstanceCasted models.CardInstance = opponentCardInstance
	interceptorContext := models.NewInterceptorContext(
		playerState,
		opponentState,
		opponentState,
		&cardInstance.Card.ID,
		&cardInstance.CardInstanceID,
		lane,
		opponentLane,
		opponentCardInstanceCasted,
	)
	err = interceptors.ExecuteInterceptors(enums.InterceptorPointHitCardBefore, &interceptorContext)
	if err != nil {
		return err
	}

	hitCard(playerState, opponentState, cardInstance, opponentCardInstance, lane, opponentLane)

	err = interceptors.ExecuteInterceptors(enums.InterceptorPointHitCardAfter, &interceptorContext)
	if err != nil {
		return err
	}

	return nil
}
