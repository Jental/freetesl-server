package models

import (
	"log"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jental/freetesl-server/common"
	"github.com/jental/freetesl-server/models/enums"
)

type PlayerMatchState struct {
	PlayerID    int
	deck        []*CardInstance
	hand        []*CardInstance
	discardPile []*CardInstance
	health      int
	runes       uint8
	mana        int
	maxMana     int
	leftLane    *Lane
	rightLane   *Lane

	hasRing      bool
	ringGemCount uint8
	isRingActive bool

	OpponentState           *PlayerMatchState
	MatchState              *Match
	Connection              *websocket.Conn
	WebsocketSendMtx        sync.Mutex
	PartiallyParsedMessages chan PartiallyParsedMessage
	Events                  chan enums.BackendEventType

	сardInstanceWaitingForAction *CardInstance // to show on FE something prophecy-action-select-like, later there'll be selects for 3-cards, but that'll be another field
	WaitingForUserActionChan     chan struct{} // TODO: make private and close in Dispose method (to be created)
}

func NewPlayerMatchState(
	playerID int,
	health int,
	runes byte,
	mana int,
	maxMana int,
	hasRing bool,
	ringGemCount uint8,
	deck []*CardInstance,
	hand []*CardInstance,
	connection *websocket.Conn,
) *PlayerMatchState {
	leftLane := NewLane(enums.LanePositionLeft, enums.LaneTypeNormal)
	rightLane := NewLane(enums.LanePositionRight, enums.LaneTypeCover)
	playerState := PlayerMatchState{
		PlayerID:    playerID,
		deck:        deck,
		hand:        hand,
		leftLane:    leftLane,
		rightLane:   rightLane,
		discardPile: make([]*CardInstance, 0),
		health:      health,
		runes:       runes,
		mana:        mana,
		maxMana:     maxMana,

		hasRing:      hasRing,
		ringGemCount: ringGemCount,
		isRingActive: hasRing,

		OpponentState:            nil,
		MatchState:               nil,
		Connection:               connection,
		PartiallyParsedMessages:  make(chan PartiallyParsedMessage, 1),
		Events:                   make(chan enums.BackendEventType, 10),
		WaitingForUserActionChan: make(chan struct{}),
	}
	leftLane.playerState = &playerState
	rightLane.playerState = &playerState
	return &playerState
}

func (playerState *PlayerMatchState) SendEvent(event enums.BackendEventType) {
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

func (playerState *PlayerMatchState) GetLane(laneID enums.LanePosition) *Lane {
	switch laneID {
	case enums.LanePositionLeft:
		return playerState.leftLane
	case enums.LanePositionRight:
		return playerState.rightLane
	}

	return nil
}

func (playerState *PlayerMatchState) GetLeftLaneCards() []*CardInstance {
	return playerState.leftLane.cardInstances
}

func (playerState *PlayerMatchState) GetRightLaneCards() []*CardInstance {
	return playerState.rightLane.cardInstances
}

func (playerState *PlayerMatchState) GetLaneCards(laneID enums.LanePosition) []*CardInstance {
	if laneID == enums.LanePositionLeft {
		return playerState.leftLane.cardInstances
	} else {
		return playerState.rightLane.cardInstances
	}
}

func (playerState *PlayerMatchState) GetAllLaneCardInstances() []*CardInstance {
	var result []*CardInstance
	result = append(result, playerState.GetLeftLaneCards()...)
	result = append(result, playerState.GetRightLaneCards()...)
	return result
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

func (playerState *PlayerMatchState) HasRing() bool { return playerState.hasRing }

func (playerState *PlayerMatchState) GetRingGemCount() uint8 { return playerState.ringGemCount }
func (playerState *PlayerMatchState) SetRingGemCount(ringGemCount uint8) {
	playerState.ringGemCount = ringGemCount

	playerState.SendEvent(enums.BackendEventRingChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentRingChanged)
}

func (playerState *PlayerMatchState) IsRingActive() bool { return playerState.isRingActive }
func (playerState *PlayerMatchState) SetRingActivity(isActive bool) {
	playerState.isRingActive = isActive

	playerState.SendEvent(enums.BackendEventRingChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentRingChanged)
}

func (playerState *PlayerMatchState) GetCardInstanceWaitingForAction() *CardInstance {
	return playerState.сardInstanceWaitingForAction
}
func (playerState *PlayerMatchState) SetCardInstanceWaitingForAction(cardInstance *CardInstance) {
	playerState.сardInstanceWaitingForAction = cardInstance

	playerState.SendEvent(enums.BackendEventCardWatingForActionChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentCardWatingForActionChanged)
}
func (playerState *PlayerMatchState) WaitForCardInstanceAction(
	onSuccess func() error,
	onTimeout func() error,
) error {
	var timeoutHappened bool
	select {
	case <-playerState.WaitingForUserActionChan:
		log.Printf("[%d]: WaitForCardInstanceAction: received card instance action completed signal", playerState.PlayerID)
		timeoutHappened = false
	case <-time.After(common.USER_ACTION_TIMEOUT * time.Second):
		log.Printf("[%d]: WaitForCardInstanceAction: user action timeout", playerState.PlayerID)
		timeoutHappened = true
	}

	playerState.сardInstanceWaitingForAction = nil

	playerState.SendEvent(enums.BackendEventCardWatingForActionChanged)
	playerState.OpponentState.SendEvent(enums.BackendEventOpponentCardWatingForActionChanged)

	if timeoutHappened {
		return onTimeout()
	} else {
		return onSuccess()
	}
}

func (playerState *PlayerMatchState) GetCardInstanceFromHand(cardInstanceID uuid.UUID) (*CardInstance, int, bool) {
	var idx = slices.IndexFunc(playerState.GetHand(), func(el *CardInstance) bool { return el.CardInstanceID == cardInstanceID })
	if idx < 0 {
		return nil, -1, false
	} else {
		return playerState.GetHand()[idx], idx, true
	}
}

func (playerState *PlayerMatchState) GetCardInstanceFromLanes(cardInstanceID uuid.UUID) (*CardInstance, *Lane, int, bool) {
	cardInstance, idx, exists := playerState.leftLane.GetCardInstance(cardInstanceID)
	if exists {
		return cardInstance, playerState.leftLane, idx, true
	}

	cardInstance, idx, exists = playerState.rightLane.GetCardInstance(cardInstanceID)
	if exists {
		return cardInstance, playerState.rightLane, idx, true
	}

	return nil, nil, -1, false
}
