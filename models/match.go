package models

import (
	"github.com/google/uuid"

	"github.com/jental/freetesl-server/common"
)

type Match struct {
	Id               uuid.UUID
	Player0State     common.Maybe[PlayerMatchState2]
	Player1State     common.Maybe[PlayerMatchState2]
	PlayerWithTurnID int
}
