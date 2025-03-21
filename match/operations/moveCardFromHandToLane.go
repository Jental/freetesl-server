package operations

import (
	"errors"
	"fmt"
	"slices"

	"github.com/jental/freetesl-server/common"
	dbEnums "github.com/jental/freetesl-server/db/enums"
	"github.com/jental/freetesl-server/match/interceptors"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

// TODO: implement with inteface after signature becomes clear
func moveCardFromHandToLaneFaceCheck(playerState *models.PlayerMatchState, cardInstance *models.CardInstance, lane *models.Lane) error {
	if !cardInstance.IsActive {
		return fmt.Errorf("[%d]: card with id '%s' is not active", playerState.PlayerID, cardInstance.CardInstanceID.String())
	}

	if cardInstance.Card.Type != dbEnums.CardTypeCreature { // TODO: there will items with mobilize later
		return fmt.Errorf("[%d]: card with id '%s' has '%d' type and cannot be moved to lane", playerState.PlayerID, cardInstance.CardInstanceID.String(), byte(cardInstance.Card.Type))
	}

	var currentMana = playerState.GetMana()
	if cardInstance.Cost > currentMana {
		return fmt.Errorf("not enough mana '%d' of '%d'", cardInstance.Cost, currentMana)
	}

	if lane.CountCardInstances() >= common.MAX_LANE_CARDS {
		return errors.New("lane is already full")
	}

	return nil
}

// logic itself
func moveCardFromHandToLane(playerState *models.PlayerMatchState, matchState *models.Match, cardInstance *models.CardInstance, cardInHandIdx int, lane *models.Lane) {
	lane.AddCardInstance(cardInstance)
	playerState.SetHand(slices.Delete(playerState.GetHand(), cardInHandIdx, cardInHandIdx+1))

	effectsWereUpdated := false

	cardInstance.IsActive = false
	if lane.Type == enums.LaneTypeCover {
		cardInstance.Effects = append(cardInstance.Effects, &models.Effect{EffectType: enums.EffectTypeCover, StartTurnID: matchState.TurnID})
		effectsWereUpdated = true
	}

	var currentMana = playerState.GetMana()
	playerState.SetMana(currentMana - cardInstance.Cost)

	if effectsWereUpdated {
		playerState.SendEvent(enums.BackendEventCardInstancesChanged)
		playerState.OpponentState.SendEvent(enums.BackendEventOpponentCardInstancesChanged)
		playerState.SendEvent(enums.BackendEventLanesChanged)
		playerState.OpponentState.SendEvent(enums.BackendEventOpponentLanesChanged)
	}
}

func MoveCardFromHandToLane(playerState *models.PlayerMatchState, matchState *models.Match, cardInstance *models.CardInstance, cardInHandIdx int, lane *models.Lane) error {
	err := moveCardFromHandToLaneFaceCheck(playerState, cardInstance, lane)
	if err != nil {
		return err
	}

	interceptorContext := models.NewInterceptorContext(
		playerState,
		nil,
		playerState,
		cardInstance.Card.ID,
		cardInstance.CardInstanceID,
		nil,
		lane,
		nil,
	)
	err = interceptors.ExecuteInterceptors(enums.InterceptorPointMoveCardFromHandToLaneBefore, &interceptorContext)
	if err != nil {
		return err
	}

	moveCardFromHandToLane(playerState, matchState, cardInstance, cardInHandIdx, lane)

	err = interceptors.ExecuteInterceptors(enums.InterceptorPointMoveCardFromHandToLaneAfter, &interceptorContext)
	if err != nil {
		return err
	}

	err = interceptors.ExecuteInterceptors(enums.InterceptorPointCardPlay, &interceptorContext)
	if err != nil {
		return err
	}

	return nil
}
