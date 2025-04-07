package operations

import (
	"fmt"
	"slices"

	"github.com/jental/freetesl-server/models"
)

// TODO: implement with inteface after signature becomes clear
func moveCardFromHandToLaneCheck(playerState *models.PlayerMatchState, cardInstance *models.CardInstance, lane *models.Lane) error {
	err := moveCardToLaneCheck(playerState, cardInstance, lane)
	if err != nil {
		return err
	}

	var currentMana = playerState.GetMana()
	if cardInstance.Cost > currentMana {
		return fmt.Errorf("not enough mana '%d' of '%d'", cardInstance.Cost, currentMana)
	}

	return nil
}

// logic itself
func moveCardFromHandToLane(playerState *models.PlayerMatchState, matchState *models.Match, cardInstance *models.CardInstance, lane *models.Lane) {
	MoveCardToLane(playerState, matchState, cardInstance, lane)

	cardInHandIdx := slices.Index(playerState.GetHand(), cardInstance)
	playerState.SetHand(slices.Delete(playerState.GetHand(), cardInHandIdx, cardInHandIdx+1))

	var currentMana = playerState.GetMana()
	playerState.SetMana(currentMana - cardInstance.Cost)
}

func MoveCardFromHandToLane(playerState *models.PlayerMatchState, matchState *models.Match, cardInstance *models.CardInstance, lane *models.Lane) error {
	err := moveCardFromHandToLaneCheck(playerState, cardInstance, lane)
	if err != nil {
		return err
	}

	moveCardFromHandToLane(playerState, matchState, cardInstance, lane)

	return nil
}
