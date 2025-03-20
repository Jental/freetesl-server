package models

import (
	"github.com/google/uuid"
)

type InterceptorContext struct {
	PlayerState        *PlayerMatchState
	OpponentState      *PlayerMatchState
	TargetPlayerState  *PlayerMatchState
	CardID             int
	CardInstanceID     uuid.UUID
	SourceLane         *Lane         // pointer - to make nullable
	TargetLane         *Lane         // pointer - to make nullable
	TargetCardInstance *CardInstance // pointer - to make nullable
}

func NewInterceptorContext(
	playerState *PlayerMatchState,
	opponentState *PlayerMatchState,
	targetPlayerState *PlayerMatchState,
	cardID int,
	cardInstanceID uuid.UUID,
	sourceLane *Lane,
	targetLane *Lane,
	targetCardInstance *CardInstance,
) InterceptorContext {
	return InterceptorContext{
		PlayerState:        playerState,
		OpponentState:      opponentState,
		CardID:             cardID,
		CardInstanceID:     cardInstanceID,
		SourceLane:         sourceLane,
		TargetPlayerState:  targetPlayerState,
		TargetLane:         targetLane,
		TargetCardInstance: targetCardInstance,
	}
}
