package mappers

import (
	"github.com/jental/freetesl-server/db/enums"
	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/models"
	"github.com/samber/lo"
)

func mapToCardDTO(model *dbModels.Card) dtos.CardDTO {
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

func MapToAllCardsDTO(model []*dbModels.Card) []*dtos.CardDTO {
	return lo.Map[*dbModels.Card, *dtos.CardDTO](
		model,
		func(item *dbModels.Card, _ int) *dtos.CardDTO { var dto = mapToCardDTO(item); return &dto })
}

func mapToPlayerInformationDTO(model *models.Player) dtos.PlayerInformationDTO {
	return dtos.PlayerInformationDTO{
		ID:         model.ID,
		Name:       model.DisplayName,
		AvatarName: model.AvatarName,
		State:      byte(model.State),
	}
}

func MapToPlayerInformationDTOs(model []*models.Player) []*dtos.PlayerInformationDTO {
	return lo.Map(model, func(item *models.Player, _ int) *dtos.PlayerInformationDTO {
		var dto = mapToPlayerInformationDTO(item)
		return &dto
	})
}
