package operations

import (
	"log"

	"github.com/jental/freetesl-server/common"
	dbEnums "github.com/jental/freetesl-server/db/enums"
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
) (bool, error) {
	// let's reduce player state from 24 to 10
	// reduce player health till next rune (e.g. from 24 to 21)
	// break rune and reduce health to 20 - to keep runes and health consistent
	// at this moment following things can happen
	// - draw card on rune break
	// - play prophecy card on rune break
	//   then it's summon effect should be triggered - it can be anything, including health and runes restore or damage
	//   - if health is restored next health reduce will be from this point (e.g from 25 to 16 instead of from 20 to 16)
	//   - if additional damage is done - it's added to total amount of health to be reduced
	//   but it's not player immediately - user need to select lane/target card (with some timeout)
	// reduce player health till next rune (e.g. from 20 to 16)
	// beak rune
	// ...
	// ? can we do it with interceptors
	//   I think, no
	// ? can rune break be aborted
	//   I don't remember such cards

	log.Printf("[%d]: ReducePlayerHealth: amount: %d", playerState.PlayerID, amount)
	amountLeft := amount

	for {
		log.Printf("[%d]: ReducePlayerHealth: amount left: %d; health: %d; runes: %d", playerState.PlayerID, amountLeft, playerState.GetHealth(), playerState.GetRunes())
		if amountLeft <= 0 {
			break
		}

		nextRuneHealth := int(playerState.GetRunes() * common.HEALTH_BETWEEN_RUNES)
		runeIsBroken := playerState.GetRunes() > 0 && amountLeft >= playerState.GetHealth()-nextRuneHealth
		var healthToSubstract int
		if runeIsBroken {
			healthToSubstract = playerState.GetHealth() - nextRuneHealth + 1 // +1 because we will will substract additional 1 health when we handle rune break
		} else {
			healthToSubstract = amountLeft
		}

		coreOperations.ReducePlayerHealth(playerState, healthToSubstract)
		amountLeft = amountLeft - healthToSubstract
		log.Printf("[%d]: ReducePlayerHealth: reduced player health: %d; amount left: %d", playerState.PlayerID, playerState.GetHealth(), amountLeft)

		if playerState.GetHealth() <= 0 {
			// TODO: there'll be an exception with a Vivec card in play later
			match.EndMatch(playerState.MatchState, playerState.OpponentState.PlayerID)
			return true, nil
		}

		if runeIsBroken {
			log.Printf("[%d]: ReducePlayerHealth: rune is broken", playerState.PlayerID)

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
				return false, err
			}

			coreOperations.ReducePlayerHealth(playerState, 1) // rune count is reduced inside
			amountLeft = amountLeft - 1
			log.Printf("[%d]: ReducePlayerHealth: reduced player health and broken rune: %d; amount left: %d", playerState.PlayerID, playerState.GetHealth(), amountLeft)

			deck := playerState.GetDeck()
			if len(deck) > 0 {
				topCard := deck[0]
				if topCard.HasKeyword(dbEnums.CardKeywordProphecy) {
					log.Printf("[%d]: ReducePlayerHealth: prophecy", playerState.PlayerID)
					topCard.IsActive = true
					playerState.SetCardInstanceWaitingForAction(topCard)
					// 0. user will select an action
					// 1. drawCard or moveCardFromDeckToLane request will be called
					// 2. and after FE will send actionCompleted request, that will write to this channel
					playerState.WaitForCardInstanceAction(
						func() error { return nil },
						func() error {
							coreOperations.DrawCard(playerState)
							return nil
						},
					)
				} else {
					log.Printf("[%d]: ReducePlayerHealth: drawing a card", playerState.PlayerID)
					coreOperations.DrawCard(playerState)
				}
			}

			err = interceptors.ExecuteInterceptors(enums.InterceptorPointRuneBreakAfter, &interceptorContext)
			if err != nil {
				return false, err
			}

			log.Printf("[%d]: ReducePlayerHealth: rune beak handling is done", playerState.PlayerID)
		}
	}

	return false, nil
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

	matchHasEnded, err := reducePlayerHealth(playerState, amount)
	if err != nil {
		return err
	}
	if matchHasEnded {
		return nil
	}

	interceptorContext = models.NewInterceptorContext(
		playerState,
		nil,
		playerState,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	err = interceptors.ExecuteInterceptors(enums.InterceptorPointHealthReduceAfter, &interceptorContext)
	if err != nil {
		return err
	}

	return nil
}
