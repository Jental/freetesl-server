package models

import (
	"time"

	"github.com/jental/freetesl-server/models/enums"
)

type PlayerRuntimeInformation struct {
	State            enums.PlayerState
	LastActivityTime *time.Time
	SelectedDeckID   *int // nullable; for LookingForOpponent state
}
