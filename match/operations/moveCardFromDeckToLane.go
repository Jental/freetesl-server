package operations

import (
	"fmt"

	"github.com/jental/freetesl-server/match/models"
)

// TODO: implement with inteface after signature becomes clear
func moveCardFromDeckToLaneCheck(playerState *models.PlayerMatchState, lane *models.Lane) error {
	deck := playerState.GetDeck()
	if len(deck) == 0 {
		return fmt.Errorf("[%d]: MoveCardFromDeckToLane: deck is empty", playerState.PlayerID)
	}

	cardInstanceCreature, ok := deck[0].(*models.CardInstanceCreature)
	if !ok {
		return fmt.Errorf("[%d]: MoveCardFromDeckToLane: top card is not a creature", playerState.PlayerID)
	}

	return moveCardToLaneCheck(playerState, cardInstanceCreature, lane)
}

// logic itself
func moveCardFromDeckToLane(playerState *models.PlayerMatchState, matchState *models.Match, lane *models.Lane) {
	// TODO:
	// - think about returning results from BE to FE
	//   or relay on FE for checks - ? what will happen in this case
	//     let's first see, what happens is there's no FE check
	//       case: user tries to add prophecy card to a lane, that's full
	//       moveCardFromDeckToLane req is sent, BE rejects this operation
	//       in parallel FE sends waitedUserActionCompleted req
	//       FE receives resetted state
	//       card is still in deck - nothing criminal happened, just prophecy card play was not applied
	//     but with FE check of course it would be better
	//   so let's go without backend results this time
	//   and it may be even better to have FE checks - instead there may be some strange behaviour on FE rollbacks (+ they have to be implemented)
	// - add ane more interceptor point - for card draw, use both (or event 3: moveToLane, draw and cardPlay) in this file

	deck := playerState.GetDeck()
	cardInstanceCreature := deck[0].(*models.CardInstanceCreature) // check has been done in the ...Check method

	MoveCardToLane(playerState, matchState, cardInstanceCreature, lane)
	playerState.SetDeck(deck[1:])
}

func MoveCardFromDeckToLane(playerState *models.PlayerMatchState, matchState *models.Match, lane *models.Lane) error {
	err := moveCardFromDeckToLaneCheck(playerState, lane)
	if err != nil {
		return err
	}

	moveCardFromDeckToLane(playerState, matchState, lane)

	// I think there should be no interceptor for card draw as it is not a real draw, but moving a card from deck to lane

	return nil
}
