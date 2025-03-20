package models

import "github.com/jental/freetesl-server/models/enums"

type Effect struct {
	EffectType  enums.EffectType
	StartTurnID int
}
