package operations

import (
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/coreOperations"
	"github.com/jental/freetesl-server/match/interceptors"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

// TODO: implement with inteface after signature becomes clear
func reducePlayerHealthCheck() error {
	return nil
}

// logic itself
func reducePlayerHealth(
	playerState *models.PlayerMatchState,
	amount int,
) error {
	originalRuneCount := playerState.GetRunes()
	coreOperations.ReducePlayerHealth(playerState, amount)
	updatedRuneCount := playerState.GetRunes()

	runeDiff := int(max(0, originalRuneCount-updatedRuneCount))
	if runeDiff != 0 {
		for i := 0; i < runeDiff; i++ {
			interceptorContext := models.NewInterceptorContext(
				playerState,
				nil,
				playerState,
				nil,
				nil,
				nil,
				nil,
				nil,
			)
			err := interceptors.ExecuteInterceptors(enums.InterceptorPointRuneBreakBefore, &interceptorContext)
			if err != nil {
				return err
			}

			if len(playerState.GetDeck()) > 0 {
				coreOperations.DrawCard(playerState)
				// TODO: trigger prophecies
			}

			err = interceptors.ExecuteInterceptors(enums.InterceptorPointRuneBreakAfter, &interceptorContext)
			if err != nil {
				return err
			}
		}
	}

	if playerState.GetHealth() <= 0 {
		// TODO: there'll be an exception with a Vivec card in play later
		match.EndMatch(playerState.MatchState, playerState.OpponentState.PlayerID)
	}

	return nil
}

func ReducePlayerHealth(
	playerState *models.PlayerMatchState,
	amount int,
) error {
	err := reducePlayerHealthCheck()
	if err != nil {
		return err
	}

	interceptorContext := models.NewInterceptorContext(
		playerState,
		nil,
		playerState,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	err = interceptors.ExecuteInterceptors(enums.InterceptorPointHealthReduceBefore, &interceptorContext)
	if err != nil {
		return err
	}

	err = reducePlayerHealth(playerState, amount)
	if err != nil {
		return err
	}

	err = interceptors.ExecuteInterceptors(enums.InterceptorPointHealthReduceAfter, &interceptorContext)
	if err != nil {
		return err
	}

	return nil
}
