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
func playItemCardCheck(playerState *models.PlayerMatchState, cardInstance *models.CardInstanceItem) error {
	if cardInstance.Card.Type != dbEnums.CardTypeItem {
		return fmt.Errorf("[%d]: PlayItemCard: card with id '%s' is not an item", playerState.PlayerID, cardInstance.CardInstanceID.String())
	}

	if !cardInstance.IsActive() {
		return fmt.Errorf("[%d]: PlayItemCard: card with id '%s' is not active", playerState.PlayerID, cardInstance.CardInstanceID.String())
	}

	var currentMana = playerState.GetMana()
	if cardInstance.Cost > currentMana {
		return fmt.Errorf("[%d]: PlayItemCard: not enough mana '%d' of '%d'", playerState.PlayerID, cardInstance.Cost, currentMana)
	}

	return nil
}

// logic itself
func playItemCard(playerState *models.PlayerMatchState, cardInstance *models.CardInstanceItem, targetCardInstance *models.CardInstanceCreature) {
	coreOperations.DiscardCardFromHand(playerState, cardInstance)
	coreOperations.AddItem(playerState, targetCardInstance, cardInstance)

	cardInstance.SetIsActive(false)
	var currentMana = playerState.GetMana()
	playerState.SetMana(currentMana - cardInstance.Cost)
}

func PlayItemCard(
	playerState *models.PlayerMatchState,
	opponentState *models.PlayerMatchState,
	cardInstance *models.CardInstanceItem,
	targetCardInstance *models.CardInstanceCreature,
	isTargetCardFromOpponent bool,
) error {
	err := playItemCardCheck(playerState, cardInstance)
	if err != nil {
		return err
	}

	playItemCard(playerState, cardInstance, targetCardInstance)

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
		nil,
		targetCardInstance,
	)
	err = interceptors.ExecuteInterceptors(enums.InterceptorPointCardPlay, &interceptorContext)
	if err != nil {
		return err
	}

	return nil
}
