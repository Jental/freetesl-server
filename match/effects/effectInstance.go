package effects

import (
	"github.com/google/uuid"
)

type EffectInstance struct {
	Effect               IEffect
	StartTurnID          int
	SourceCardInstanceID *uuid.UUID // nullable
}

func NewEffectInstance(effect IEffect, startTurnID int, sourceCardInstanceID *uuid.UUID) EffectInstance {
	return EffectInstance{
		Effect:               effect,
		StartTurnID:          startTurnID,
		SourceCardInstanceID: sourceCardInstanceID,
	}
}
