package models

import (
	"github.com/google/uuid"
	dbEnums "github.com/jental/freetesl-server/db/enums"
)

type KeywordInstance struct {
	Keyword              dbEnums.CardKeyword
	StartTurnID          *int      // nullable
	SourceCardInstanceID uuid.UUID // nullable
}

func NewKeywordInstance(
	keyword dbEnums.CardKeyword,
	startTurnID *int,
	sourceCardInstaneID uuid.UUID,
) KeywordInstance {
	return KeywordInstance{
		Keyword:              keyword,
		StartTurnID:          startTurnID,
		SourceCardInstanceID: sourceCardInstaneID,
	}
}
