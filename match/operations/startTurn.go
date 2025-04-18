package operations

import (
	"slices"

	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/models/enums"
)

// TODO: implement with inteface after signature becomes clear
func startTurnCheck() error {
	return nil
}

// logic itself
func startTurn(playerState *models.PlayerMatchState, matchState *models.Match) {
	coreOperations.StartTurn(playerState)
	playerState.OpponentState.SetRingActivity(false)

	effectsWereUpdated := false
	for _, card := range playerState.GetAllLaneCardInstances() {
		if card.HasEffect(enums.EffectTypeShackled) {
			originalLen := len(card.Effects)
			card.Effects = slices.DeleteFunc(card.Effects, func(eff *models.Effect) bool {
				return eff.EffectType == enums.EffectTypeShackled && matchState.TurnID-eff.StartTurnID > common.SHACKLE_TURNS_TO_SKIP
			})
			effectsWereUpdated = effectsWereUpdated || originalLen != len(card.Effects)
		}
		if card.HasEffect(enums.EffectTypeCover) {
			originalLen := len(card.Effects)
			card.Effects = slices.DeleteFunc(card.Effects, func(eff *models.Effect) bool {
				return eff.EffectType == enums.EffectTypeCover
			})
			effectsWereUpdated = effectsWereUpdated || originalLen != len(card.Effects)
		}

		card.SetIsActive(!card.HasEffect(enums.EffectTypeShackled))
	}

	if effectsWereUpdated {
		playerState.SendEvent(enums.BackendEventCardInstancesChanged)
		playerState.OpponentState.SendEvent(enums.BackendEventOpponentCardInstancesChanged)
		playerState.SendEvent(enums.BackendEventLanesChanged)
		playerState.OpponentState.SendEvent(enums.BackendEventOpponentLanesChanged)
	}
}

func StartTurn(playerState *models.PlayerMatchState, matchState *models.Match) {
	_ = startTurnCheck()
	startTurn(playerState, matchState)
}
