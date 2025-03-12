package models

import (
	"log"
	"sync"

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

	OpponentState           *PlayerMatchState
	MatchState              *Match
	Connection              *websocket.Conn
	WebsocketSendMtx        sync.Mutex
	PartiallyParsedMessages chan PartiallyParsedMessage
	Events                  chan enums.BackendEventType
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

		OpponentState:           nil,
		MatchState:              nil,
		Connection:              connection,
		PartiallyParsedMessages: make(chan PartiallyParsedMessage, 1),
		Events:                  make(chan enums.BackendEventType, 10),
	}
}

func (playerState *PlayerMatchState) SendEvent(event enums.BackendEventType) {
	// TODO: use it everywhere instead of sending to Events channel
	if playerState != nil && playerState.Connection != nil {
		playerState.Events <- event
	} else {
		// TODO: probably add error return
		log.Printf("[%d]: writing to Events channel that should be closed", playerState.PlayerID)
	}
}

func (playerState *PlayerMatchState) GetDeck() []*CardInstance { return playerState.deck }
func (playerState *PlayerMatchState) SetDeck(deck []*CardInstance) {
	playerState.deck = deck

	playerState.SendEvent(enums.BackendEventDeckChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentDeckChanged)
}

func (playerState *PlayerMatchState) GetHand() []*CardInstance { return playerState.hand }
func (playerState *PlayerMatchState) SetHand(hand []*CardInstance) {
	playerState.hand = hand

	playerState.SendEvent(enums.BackendEventHandChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentHandChanged)
}

func (playerState *PlayerMatchState) GetDiscardPile() []*CardInstance { return playerState.discardPile }
func (playerState *PlayerMatchState) SetDiscardPile(discardPile []*CardInstance) {
	playerState.discardPile = discardPile

	playerState.SendEvent(enums.BackendEventDiscardPileChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentDiscardPileChanged)
}

func (playerState *PlayerMatchState) GetLeftLaneCards() []*CardInstance {
	return playerState.leftLaneCards
}
func (playerState *PlayerMatchState) SetLeftLaneCards(leftLaneCards []*CardInstance) {
	playerState.leftLaneCards = leftLaneCards

	playerState.SendEvent(enums.BackendEventLanesChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentLanesChanged)
}

func (playerState *PlayerMatchState) GetRightLaneCards() []*CardInstance {
	return playerState.rightLaneCards
}
func (playerState *PlayerMatchState) SetRightLaneCards(rightLaneCards []*CardInstance) {
	playerState.rightLaneCards = rightLaneCards

	playerState.SendEvent(enums.BackendEventLanesChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentLanesChanged)
}

func (playerState *PlayerMatchState) GetLaneCards(laneID enums.Lane) []*CardInstance {
	if laneID == enums.LaneRight {
		return playerState.GetRightLaneCards()
	} else {
		return playerState.GetLeftLaneCards()
	}
}

func (playerState *PlayerMatchState) GetHealth() int { return playerState.health }
func (playerState *PlayerMatchState) SetHealth(health int) {
	playerState.health = health

	playerState.SendEvent(enums.BackendEventHealthChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentHealthChanged)
}

func (playerState *PlayerMatchState) GetRunes() byte { return playerState.runes }
func (playerState *PlayerMatchState) SetRunes(runes byte) {
	playerState.runes = runes

	playerState.SendEvent(enums.BackendEventHealthChanged) // decided to reuse event
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentHealthChanged)
}

func (playerState *PlayerMatchState) GetMana() int { return playerState.mana }
func (playerState *PlayerMatchState) SetMana(mana int) {
	playerState.mana = mana

	playerState.SendEvent(enums.BackendEventManaChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentManaChanged)
}

func (playerState *PlayerMatchState) GetMaxMana() int { return playerState.maxMana }
func (playerState *PlayerMatchState) SetMaxMana(maxMana int) {
	playerState.maxMana = maxMana

	playerState.SendEvent(enums.BackendEventManaChanged) // decided to reuse event
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentManaChanged)
}
