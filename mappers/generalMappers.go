package mappers

import (
	"github.com/jental/freetesl-server/db/enums"
	"github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/dtos"
	"github.com/samber/lo"
)

func mapToCardDTO(model *models.Card) dtos.CardDTO {
	return dtos.CardDTO{
		ID:     model.ID,
		Type:   byte(model.Type),
		Power:  model.Power,
		Health: model.Health,
		Cost:   model.Cost,
		Keywords: lo.Map[enums.CardKeyword, int](
			model.Keywords,
			func(item enums.CardKeyword, _ int) int { return int(item) }),
	}
}

func MapToAllCardsDTO(model []*models.Card) []*dtos.CardDTO {
	return lo.Map[*models.Card, *dtos.CardDTO](
		model,
		func(item *models.Card, _ int) *dtos.CardDTO { var dto = mapToCardDTO(item); return &dto })
}
