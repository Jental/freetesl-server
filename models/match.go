package models

import (
	"github.com/google/uuid"

	"github.com/jental/freetesl-server/common"
)

type Match struct {
	Id               uuid.UUID
	Player0State     common.Maybe[PlayerMatchState]
	Player1State     common.Maybe[PlayerMatchState]
	PlayerWithTurnID int
	WinnerID         int
}
