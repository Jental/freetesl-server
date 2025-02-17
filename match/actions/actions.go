package actions

import (
	"errors"
	"fmt"
	"slices"

	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/match/senders"
	"github.com/jental/freetesl-server/models"
)

func DrawCard(playerState *models.PlayerMatchState2) {
	if len(playerState.Deck) == 0 {
		return // TODO: for now doing nothing, but later next rune should be broken
	}

	var drawnCard = playerState.Deck[0]
	drawnCard.IsActive = true
	playerState.Hand = append(playerState.Hand, drawnCard)
	playerState.Deck = playerState.Deck[1:]
}

func SwitchTurn(match *models.Match) {
	var isFirstPlayersTurn = match.Player0State.Value.PlayerID == match.PlayerWithTurnID
	if isFirstPlayersTurn {
		match.PlayerWithTurnID = match.Player1State.Value.PlayerID
	} else {
		match.PlayerWithTurnID = match.Player0State.Value.PlayerID
	}
}

func MoveCardToLane(playerState *models.PlayerMatchState2, cardInstance *models.CardInstance, cardInHandIdx int, laneID byte) error {
	if cardInstance.Cost > playerState.Mana {
		return fmt.Errorf("not enough mana '%d' of '%d'", cardInstance.Cost, playerState.Mana)
	}

	if laneID == common.LEFT_LANE_ID {
		if len(playerState.LeftLaneCards) >= common.MAX_LANE_CARDS {
			return errors.New("lane is already full")
		}
		playerState.LeftLaneCards = append(playerState.LeftLaneCards, cardInstance)
	} else if laneID == common.RIGHT_LANE_ID {
		if len(playerState.RightLaneCards) >= common.MAX_LANE_CARDS {
			return errors.New("lane is already full")
		}
		playerState.RightLaneCards = append(playerState.RightLaneCards, cardInstance)
	} else {
		return fmt.Errorf("invalid lane id: %d", laneID)
	}

	cardInstance.IsActive = false
	playerState.Hand = slices.Delete(playerState.Hand, cardInHandIdx, cardInHandIdx+1)
	playerState.Mana = playerState.Mana - cardInstance.Cost

	return nil
}

func ReducePlayerHealth(playerState *models.PlayerMatchState2, matchState *models.Match, amount int) {
	playerState.Health = playerState.Health - amount

	var expectedRuneCount uint8 = uint8((playerState.Health - 1) / 5)
	var runeCount = max(0, min(expectedRuneCount, playerState.Runes))
	playerState.Runes = runeCount

	// TODO: trigger prophecies

	if playerState.Health <= 0 {
		senders.SendMatchEndToEveryone(matchState)
		// TODO: there'll be an exception with a Vivec card in play later
	}
}

func ReduceCardHealth(cardInstance *models.CardInstance, laneID byte, playerState *models.PlayerMatchState2, matchState *models.Match, amount int) error {
	cardInstance.Health = cardInstance.Health - amount

	if cardInstance.Health <= 0 {
		if laneID == common.LEFT_LANE_ID {
			var idx = slices.Index(playerState.LeftLaneCards, cardInstance)
			if idx < 0 {
				return errors.New("player does have the card in a left lane.")
			}

			playerState.LeftLaneCards = slices.Delete(playerState.LeftLaneCards, idx, idx+1)
			playerState.DiscardPile = append(playerState.DiscardPile, cardInstance)
			// senders.SendDiscardPileToEveryone(matchState) - I don't like it
			// TODO: trigger some event, based on which it will be decided which updates to send
			return nil

		} else if laneID == common.RIGHT_LANE_ID {
			var idx = slices.Index(playerState.RightLaneCards, cardInstance)
			if idx < 0 {
				return errors.New("player does have the card in a right lane.")
			}

			playerState.RightLaneCards = slices.Delete(playerState.RightLaneCards, idx, idx+1)
			playerState.DiscardPile = append(playerState.DiscardPile, cardInstance)
			return nil

		} else {
			return fmt.Errorf("invalid lane ID: %d", laneID)
		}
	}

	return nil
}
