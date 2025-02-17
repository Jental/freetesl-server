package models

import (
	"github.com/gorilla/websocket"
	"github.com/jental/freetesl-server/models/enums"
)

type PlayerMatchState struct {
	PlayerID       int
	deck           []*CardInstance
	hand           []*CardInstance
	discardPile    []*CardInstance
	health         int
	runes          uint8
	mana           int
	maxMana        int
	leftLaneCards  []*CardInstance
	rightLaneCards []*CardInstance

	OpponentState *PlayerMatchState
	MatchState    *Match
	Connection    *websocket.Conn
	Events        chan enums.BackendEventType
}

func NewPlayerMatchState(
	playerID int,
	health int,
	runes byte,
	mana int,
	maxMana int,
	deck []*CardInstance,
	hand []*CardInstance,
	discardPile []*CardInstance,
	leftLaneCards []*CardInstance,
	rightLaneCards []*CardInstance,
	connection *websocket.Conn,
) PlayerMatchState {
	return PlayerMatchState{
		PlayerID:       playerID,
		deck:           deck,
		hand:           hand,
		leftLaneCards:  leftLaneCards,
		rightLaneCards: rightLaneCards,
		discardPile:    discardPile,
		health:         health,
		runes:          runes,
		mana:           mana,
		maxMana:        maxMana,

		OpponentState: nil,
		MatchState:    nil,
		Connection:    connection,
		Events:        make(chan enums.BackendEventType, 10),
	}
}

func (playerState *PlayerMatchState) GetDeck() []*CardInstance { return playerState.deck }
func (playerState *PlayerMatchState) SetDeck(deck []*CardInstance) {
	playerState.deck = deck

	playerState.Events <- enums.BackendEventDeckChanged
	if playerState.OpponentState != nil {
		playerState.OpponentState.Events <- enums.BackendEventOpponentDeckChanged
	}
}

func (playerState *PlayerMatchState) GetHand() []*CardInstance { return playerState.hand }
func (playerState *PlayerMatchState) SetHand(hand []*CardInstance) {
	playerState.hand = hand

	playerState.Events <- enums.BackendEventHandChanged
	if playerState.OpponentState != nil {
		playerState.OpponentState.Events <- enums.BackendEventOpponentHandChanged
	}
}

func (playerState *PlayerMatchState) GetDiscardPile() []*CardInstance { return playerState.discardPile }
func (playerState *PlayerMatchState) SetDiscardPile(discardPile []*CardInstance) {
	playerState.discardPile = discardPile

	playerState.Events <- enums.BackendEventDiscardPileChanged
	if playerState.OpponentState != nil {
		playerState.OpponentState.Events <- enums.BackendEventOpponentDiscardPileChanged
	}
}

func (playerState *PlayerMatchState) GetLeftLaneCards() []*CardInstance {
	return playerState.leftLaneCards
}
func (playerState *PlayerMatchState) SetLeftLaneCards(leftLaneCards []*CardInstance) {
	playerState.leftLaneCards = leftLaneCards

	playerState.Events <- enums.BackendEventLanesChanged
	if playerState.OpponentState != nil {
		playerState.OpponentState.Events <- enums.BackendEventOpponentLanesChanged
	}
}

func (playerState *PlayerMatchState) GetRightLaneCards() []*CardInstance {
	return playerState.rightLaneCards
}
func (playerState *PlayerMatchState) SetRightLaneCards(rightLaneCards []*CardInstance) {
	playerState.rightLaneCards = rightLaneCards

	playerState.Events <- enums.BackendEventLanesChanged
	if playerState.OpponentState != nil {
		playerState.OpponentState.Events <- enums.BackendEventOpponentLanesChanged
	}
}

func (playerState *PlayerMatchState) GetHealth() int { return playerState.health }
func (playerState *PlayerMatchState) SetHealth(health int) {
	playerState.health = health

	playerState.Events <- enums.BackendEventHealthChanged
	if playerState.OpponentState != nil {
		playerState.OpponentState.Events <- enums.BackendEventOpponentHealthChanged
	}
}

func (playerState *PlayerMatchState) GetRunes() byte { return playerState.runes }
func (playerState *PlayerMatchState) SetRunes(runes byte) {
	playerState.runes = runes

	playerState.Events <- enums.BackendEventHealthChanged // decided to reuse event
	if playerState.OpponentState != nil {
		playerState.OpponentState.Events <- enums.BackendEventOpponentHealthChanged
	}
}

func (playerState *PlayerMatchState) GetMana() int { return playerState.mana }
func (playerState *PlayerMatchState) SetMana(mana int) {
	playerState.mana = mana

	playerState.Events <- enums.BackendEventManaChanged
	if playerState.OpponentState != nil {
		playerState.OpponentState.Events <- enums.BackendEventOpponentManaChanged
	}
}

func (playerState *PlayerMatchState) GetMaxMana() int { return playerState.maxMana }
func (playerState *PlayerMatchState) SetMaxMana(maxMana int) {
	playerState.maxMana = maxMana

	playerState.Events <- enums.BackendEventManaChanged // decided to reuse event
	if playerState.OpponentState != nil {
		playerState.OpponentState.Events <- enums.BackendEventOpponentManaChanged
	}
}
