package models

import dbModels "github.com/jental/freetesl-server/db/models"

type Deck struct {
	ID         int
	Name       string
	AvatarName string
	PlayerID   int
	Cards      []*CardWithCount
	Attributes []*dbModels.Attribute
}
