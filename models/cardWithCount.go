package models

import dbModels "github.com/jental/freetesl-server/db/models"

type CardWithCount struct {
	Card  *dbModels.Card
	Count int
}
