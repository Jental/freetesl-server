package models

import (
	"github.com/google/uuid"
	"github.com/jental/freetesl-server/models/enums"
)

type InterceptorContext struct {
	PlayerState        *PlayerMatchState
	OpponentState      *PlayerMatchState
	TargetPlayerState  *PlayerMatchState
	CardID             int
	CardInstanceID     uuid.UUID
	SourceLaneID       *enums.Lane   // pointer - to make nullable
	TargetLaneID       *enums.Lane   // pointer - to make nullable
	TargetCardInstance *CardInstance // pointer - to make nullable
}

func NewInterceptorContext(
	playerState *PlayerMatchState,
	opponentState *PlayerMatchState,
	targetPlayerState *PlayerMatchState,
	cardID int,
	cardInstanceID uuid.UUID,
	sourceLaneID *enums.Lane,
	targetLaneID *enums.Lane,
	targetCardInstance *CardInstance,
) InterceptorContext {
	return InterceptorContext{
		PlayerState:        playerState,
		OpponentState:      opponentState,
		CardID:             cardID,
		CardInstanceID:     cardInstanceID,
		SourceLaneID:       sourceLaneID,
		TargetPlayerState:  targetPlayerState,
		TargetLaneID:       targetLaneID,
		TargetCardInstance: targetCardInstance,
	}
}
