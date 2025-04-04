package operations

import (
	"errors"
	"fmt"

	"github.com/jental/freetesl-server/common"
	dbEnums "github.com/jental/freetesl-server/db/enums"
	"github.com/jental/freetesl-server/match/interceptors"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

// TODO: implement with inteface after signature becomes clear
func moveCardToLaneCheck(playerState *models.PlayerMatchState, cardInstance *models.CardInstance, lane *models.Lane) error {
	if !cardInstance.IsActive {
		return fmt.Errorf("[%d]: card with id '%s' is not active", playerState.PlayerID, cardInstance.CardInstanceID.String())
	}

	if cardInstance.Card.Type != dbEnums.CardTypeCreature { // TODO: there will items with mobilize later
		return fmt.Errorf("[%d]: card with id '%s' has '%d' type and cannot be moved to lane", playerState.PlayerID, cardInstance.CardInstanceID.String(), byte(cardInstance.Card.Type))
	}

	if lane.CountCardInstances() >= common.MAX_LANE_CARDS {
		return errors.New("lane is already full")
	}

	return nil
}

// logic itself
func moveCardToLane(playerState *models.PlayerMatchState, matchState *models.Match, cardInstance *models.CardInstance, lane *models.Lane) {
	lane.AddCardInstance(cardInstance)
	cardInstance.IsActive = false

	// TODO: maybe do it through interceptor
	effectsWereUpdated := false
	if lane.Type == enums.LaneTypeCover && !cardInstance.HasKeyword(dbEnums.CardKeywordGuard) {
		cardInstance.Effects = append(cardInstance.Effects, &models.Effect{EffectType: enums.EffectTypeCover, StartTurnID: matchState.TurnID})
		effectsWereUpdated = true
	}

	if effectsWereUpdated {
		playerState.SendEvent(enums.BackendEventCardInstancesChanged)
		playerState.OpponentState.SendEvent(enums.BackendEventOpponentCardInstancesChanged)
		playerState.SendEvent(enums.BackendEventLanesChanged) // to trigger FE redraw
		playerState.OpponentState.SendEvent(enums.BackendEventOpponentLanesChanged)
	}
}

func MoveCardToLane(playerState *models.PlayerMatchState, matchState *models.Match, cardInstance *models.CardInstance, cardInHandIdx int, lane *models.Lane) error {
	err := moveCardToLaneCheck(playerState, cardInstance, lane)
	if err != nil {
		return err
	}

	interceptorContext := models.NewInterceptorContext(
		playerState,
		nil,
		playerState,
		&cardInstance.Card.ID,
		&cardInstance.CardInstanceID,
		nil,
		lane,
		nil,
	)
	err = interceptors.ExecuteInterceptors(enums.InterceptorPointMoveCardToLaneBefore, &interceptorContext)
	if err != nil {
		return err
	}

	moveCardToLane(playerState, matchState, cardInstance, lane)

	err = interceptors.ExecuteInterceptors(enums.InterceptorPointMoveCardToLaneAfter, &interceptorContext)
	if err != nil {
		return err
	}

	err = interceptors.ExecuteInterceptors(enums.InterceptorPointCardPlay, &interceptorContext)
	if err != nil {
		return err
	}

	return nil
}
