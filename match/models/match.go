package models

import (
	"github.com/google/uuid"

	"github.com/jental/freetesl-server/common"
)

type Match struct {
	Id                    uuid.UUID
	Player0State          common.Maybe[PlayerMatchState] // remove maybe
	Player1State          common.Maybe[PlayerMatchState] // remove maybe
	TurnID                int
	PlayerWithTurnID      int
	PlayerWithFirstTurnID int
	WinnerID              int // -1 means 'no winner yet' or 'draw'
}
