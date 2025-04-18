package models

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/models/enums"
)

type Lane struct {
	Position      enums.LanePosition
	Type          enums.LaneType
	cardInstances []*CardInstanceCreature
	playerState   *PlayerMatchState
}

func NewLane(position enums.LanePosition, laneType enums.LaneType) *Lane {
	return &Lane{
		Position:      position,
		Type:          laneType,
		cardInstances: make([]*CardInstanceCreature, 0),
	}
}

func (lane *Lane) GetCardInstance(id uuid.UUID) (*CardInstanceCreature, int, bool) {
	var idx = slices.IndexFunc(lane.cardInstances, func(el *CardInstanceCreature) bool { return el.CardInstanceID == id })
	if idx >= 0 {
		return lane.cardInstances[idx], idx, true
	} else {
		return nil, -1, false
	}
}

func (lane *Lane) CountCardInstances() int {
	return len(lane.cardInstances)
}

func (lane *Lane) AddCardInstance(cardInstance *CardInstanceCreature) {
	lane.cardInstances = append(lane.cardInstances, cardInstance)

	lane.playerState.SendEvent(enums.BackendEventLanesChanged)
	lane.playerState.OpponentState.SendEvent(enums.BackendEventOpponentLanesChanged)
}

func (lane *Lane) RemoveCardInstance(cardInstance *CardInstanceCreature) error {
	idx := slices.Index(lane.cardInstances, cardInstance)
	if idx < 0 {
		return fmt.Errorf("CardInstance with id '%s' is not found", cardInstance.CardInstanceID)
	}

	lane.cardInstances = slices.Delete(lane.cardInstances, idx, idx+1)

	lane.playerState.SendEvent(enums.BackendEventLanesChanged)
	lane.playerState.OpponentState.SendEvent(enums.BackendEventOpponentLanesChanged)

	return nil
}

func (lane *Lane) RemoveCardInstanceByIndex(cardInstance *CardInstanceCreature, idx int) error {
	if idx >= len(lane.cardInstances) {
		return fmt.Errorf("CardInstance index '%d' is out of range", idx)
	}
	foundCardInstance := lane.cardInstances[idx]
	if foundCardInstance.CardInstanceID != cardInstance.CardInstanceID {
		return fmt.Errorf("CardInstance by index '%d' has different id", idx)
	}
	lane.cardInstances = slices.Delete(lane.cardInstances, idx, idx+1)
	return nil
}
