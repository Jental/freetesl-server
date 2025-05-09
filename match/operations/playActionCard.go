package operations

import (
	"fmt"

	dbEnums "github.com/jental/freetesl-server/db/enums"
	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/match/interceptors"
	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/models/enums"
)

// TODO: implement with inteface after signature becomes clear
func playActionCardCheck(playerState *models.PlayerMatchState, cardInstance *models.CardInstanceAction) error {
	if cardInstance.Card.Type != dbEnums.CardTypeAction {
		return fmt.Errorf("[%d]: card with id '%s' is not an action", playerState.PlayerID, cardInstance.CardInstanceID.String())
	}

	if !cardInstance.IsActive() {
		return fmt.Errorf("[%d]: card with id '%s' is not active", playerState.PlayerID, cardInstance.CardInstanceID.String())
	}

	var currentMana = playerState.GetMana()
	if cardInstance.Cost > currentMana {
		return fmt.Errorf("not enough mana '%d' of '%d'", cardInstance.Cost, currentMana)
	}

	return nil
}

// logic itself
func playActionCard(playerState *models.PlayerMatchState, cardInstance *models.CardInstanceAction) {
	coreOperations.DiscardCardFromHand(playerState, cardInstance)
	cardInstance.SetIsActive(false)
	var currentMana = playerState.GetMana()
	playerState.SetMana(currentMana - cardInstance.Cost)
}

func PlayActionCard(
	playerState *models.PlayerMatchState,
	opponentState *models.PlayerMatchState,
	cardInstance *models.CardInstanceAction,
	targetCardInstance *models.CardInstanceCreature,
	isTargetCardFromOpponent bool,
	targetLane *models.Lane,
) error {
	err := playActionCardCheck(playerState, cardInstance)
	if err != nil {
		return err
	}

	playActionCard(playerState, cardInstance)

	targetPlayerState := opponentState
	if !isTargetCardFromOpponent {
		targetPlayerState = playerState
	}
	interceptorContext := models.NewInterceptorContext(
		playerState,
		opponentState,
		targetPlayerState,
		&cardInstance.Card.ID,
		&cardInstance.CardInstanceID,
		nil,
		targetLane,
		targetCardInstance,
	)
	err = interceptors.ExecuteInterceptors(enums.InterceptorPointCardPlay, &interceptorContext)
	if err != nil {
		return err
	}

	return nil
}
