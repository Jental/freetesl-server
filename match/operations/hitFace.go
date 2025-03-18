package operations

import (
	"fmt"

	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/match/interceptors"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

// TODO: implement with inteface after signature becomes clear
func hitFaceCheck(playerState *models.PlayerMatchState, cardInstance *models.CardInstance) error {
	if !cardInstance.IsActive {
		return fmt.Errorf("[%d]: card with id '%s' is not active", playerState.PlayerID, cardInstance.CardInstanceID.String())
	}

	return nil
}

// logic itself
func hitFace(opponentState *models.PlayerMatchState, cardInstance *models.CardInstance) {
	coreOperations.ReducePlayerHealth(opponentState, cardInstance.Power)
	cardInstance.IsActive = false
}

func HitFace(playerState *models.PlayerMatchState, opponentState *models.PlayerMatchState, cardInstance *models.CardInstance, laneID enums.Lane) error {
	err := hitFaceCheck(playerState, cardInstance)
	if err != nil {
		return err
	}

	interceptorContext := models.NewInterceptorContext(
		playerState,
		opponentState,
		opponentState,
		&cardInstance.Card.ID,
		&laneID,
		nil,
		nil,
	)
	err = interceptors.ExecuteInterceptors(enums.InterceptorPointHitFaceBefore, &interceptorContext)
	if err != nil {
		return err
	}

	hitFace(opponentState, cardInstance)

	err = interceptors.ExecuteInterceptors(enums.InterceptorPointHitFaceAfter, &interceptorContext)
	if err != nil {
		return err
	}

	return nil
}
