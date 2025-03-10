package actions

import (
	"errors"
	"fmt"
	"slices"

	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

func DrawCard(playerState *models.PlayerMatchState) {
	var deck = playerState.GetDeck()
	if len(deck) == 0 {
		return // TODO: for now doing nothing, but later next rune should be broken
	}

	var drawnCard = deck[0]
	drawnCard.IsActive = true
	playerState.SetHand(append(playerState.GetHand(), drawnCard))
	playerState.SetDeck(deck[1:])
}

func SwitchTurn(matchState *models.Match) {
	var isFirstPlayersTurn = matchState.Player0State.Value.PlayerID == matchState.PlayerWithTurnID
	if isFirstPlayersTurn {
		matchState.PlayerWithTurnID = matchState.Player1State.Value.PlayerID
	} else {
		matchState.PlayerWithTurnID = matchState.Player0State.Value.PlayerID
	}

	matchState.Player0State.Value.SendEvent(enums.BackendEventSwitchTurn)
	matchState.Player1State.Value.SendEvent(enums.BackendEventSwitchTurn)
}

func MoveCardToLane(playerState *models.PlayerMatchState, cardInstance *models.CardInstance, cardInHandIdx int, laneID byte) error {
	var currentMana = playerState.GetMana()

	if cardInstance.Cost > currentMana {
		return fmt.Errorf("not enough mana '%d' of '%d'", cardInstance.Cost, currentMana)
	}

	if laneID == common.LEFT_LANE_ID {
		if len(playerState.GetLeftLaneCards()) >= common.MAX_LANE_CARDS {
			return errors.New("lane is already full")
		}
		playerState.SetLeftLaneCards(append(playerState.GetLeftLaneCards(), cardInstance))
	} else if laneID == common.RIGHT_LANE_ID {
		if len(playerState.GetRightLaneCards()) >= common.MAX_LANE_CARDS {
			return errors.New("lane is already full")
		}
		playerState.SetRightLaneCards(append(playerState.GetRightLaneCards(), cardInstance))
	} else {
		return fmt.Errorf("invalid lane id: %d", laneID)
	}

	cardInstance.IsActive = false
	playerState.SetHand(slices.Delete(playerState.GetHand(), cardInHandIdx, cardInHandIdx+1))
	playerState.SetMana(currentMana - cardInstance.Cost)

	return nil
}

func ReducePlayerHealth(playerState *models.PlayerMatchState, amount int) {
	var updatedHealth = playerState.GetHealth() - amount
	playerState.SetHealth(updatedHealth)

	var expectedRuneCount uint8 = uint8((updatedHealth - 1) / 5)
	var runeCount = max(0, min(expectedRuneCount, playerState.GetRunes()))
	playerState.SetRunes(runeCount)

	// TODO: trigger prophecies

	if updatedHealth <= 0 {
		// TODO: there'll be an exception with a Vivec card in play later
		match.EndMatch(playerState.MatchState, playerState.OpponentState.PlayerID)
	}
}

func ReduceCardHealth(playerState *models.PlayerMatchState, cardInstance *models.CardInstance, laneID byte, amount int) error {
	cardInstance.Health = cardInstance.Health - amount

	if cardInstance.Health <= 0 {
		if laneID != common.LEFT_LANE_ID && laneID != common.RIGHT_LANE_ID {
			return fmt.Errorf("invalid lane ID: %d", laneID)
		}

		if laneID == common.LEFT_LANE_ID {
			var idx = slices.Index(playerState.GetLeftLaneCards(), cardInstance)
			if idx < 0 {
				return errors.New("player does have the card in a left lane")
			}

			playerState.SetLeftLaneCards(slices.Delete(playerState.GetLeftLaneCards(), idx, idx+1))
			playerState.SetDiscardPile(append(playerState.GetDiscardPile(), cardInstance))

		} else {
			var idx = slices.Index(playerState.GetRightLaneCards(), cardInstance)
			if idx < 0 {
				return errors.New("player does have the card in a right lane")
			}

			playerState.SetRightLaneCards(slices.Delete(playerState.GetRightLaneCards(), idx, idx+1))
			playerState.SetDiscardPile(append(playerState.GetDiscardPile(), cardInstance))
		}
	}

	playerState.SendEvent(enums.BackendEventLanesChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentLanesChanged)

	return nil
}
