package models

import "github.com/google/uuid"

type ActionContext struct {
	PlayerState        *PlayerMatchState
	OpponentState      *PlayerMatchState
	CardID             *int
	CardInstanceID     *uuid.UUID
	ParametersValues   *string
	TargetPlayerState  *PlayerMatchState // to be able to modify it's hand or lane cards
	TargetCardInstance CardInstance
	TargetLane         *Lane
}

type Action interface {
	Execute(context *ActionContext) error
}

func NewActionContext(
	playerState *PlayerMatchState,
	opponentState *PlayerMatchState,
	cardID *int,
	cardInstanceID *uuid.UUID,
	parametersValues *string,
	targetPlayerState *PlayerMatchState,
	targetCardInstance CardInstance,
	targetLane *Lane,
) ActionContext {
	return ActionContext{
		PlayerState:        playerState,
		OpponentState:      opponentState,
		CardID:             cardID,
		CardInstanceID:     cardInstanceID,
		ParametersValues:   parametersValues,
		TargetPlayerState:  targetPlayerState,
		TargetCardInstance: targetCardInstance,
		TargetLane:         targetLane,
	}
}
