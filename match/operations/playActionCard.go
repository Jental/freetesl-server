package operations

import (
	"fmt"

	dbEnums "github.com/jental/freetesl-server/db/enums"
	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/match/interceptors"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

// TODO: implement with inteface after signature becomes clear
func playActionCheck(playerState *models.PlayerMatchState, cardInstance *models.CardInstance) error {
	if cardInstance.Card.Type != dbEnums.CardTypeAction {
		return fmt.Errorf("[%d]: card with id '%s' is not an action", playerState.PlayerID, cardInstance.CardInstanceID.String())
	}

	if !cardInstance.IsActive {
		return fmt.Errorf("[%d]: card with id '%s' is not active", playerState.PlayerID, cardInstance.CardInstanceID.String())
	}

	return nil
}

// logic itself
func playAction(playerState *models.PlayerMatchState, cardInstance *models.CardInstance) {
	coreOperations.DiscardCardFromHand(playerState, cardInstance)
	cardInstance.IsActive = false
}

func PlayActionCard(
	playerState *models.PlayerMatchState,
	opponentState *models.PlayerMatchState,
	cardInstance *models.CardInstance,
	targetCardInstance *models.CardInstance,
	isTargetCardFromOpponent bool,
	targetLaneID *enums.Lane,
) error {
	err := playActionCheck(playerState, cardInstance)
	if err != nil {
		return err
	}

	playAction(playerState, cardInstance)

	targetPlayerState := opponentState
	if !isTargetCardFromOpponent {
		targetPlayerState = playerState
	}
	interceptorContext := models.NewInterceptorContext(
		playerState,
		opponentState,
		targetPlayerState,
		cardInstance.Card.ID,
		cardInstance.CardInstanceID,
		nil,
		targetLaneID,
		targetCardInstance,
	)
	err = interceptors.ExecuteInterceptors(enums.InterceptorPointCardPlay, &interceptorContext)
	if err != nil {
		return err
	}

	return nil
}
