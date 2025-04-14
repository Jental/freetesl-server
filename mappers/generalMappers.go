package mappers

import (
	"slices"

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
		Keywords: lo.Map(
			model.Keywords,
			func(item enums.CardKeyword, _ int) int { return int(item) }),
	}
}

func MapToAllCardsDTO(model []*dbModels.Card) []*dtos.CardDTO {
	return lo.Map(
		model,
		func(item *dbModels.Card, _ int) *dtos.CardDTO { var dto = mapToCardDTO(item); return &dto })
}

func MapToPlayerInformationDTO(model *models.Player) dtos.PlayerInformationDTO {
	return dtos.PlayerInformationDTO{
		ID:         model.ID,
		Name:       model.DisplayName,
		AvatarName: model.AvatarName,
		State:      byte(model.State),
	}
}

func MapToPlayerInformationDTOs(model []*models.Player) []*dtos.PlayerInformationDTO {
	return lo.Map(model, func(item *models.Player, _ int) *dtos.PlayerInformationDTO {
		var dto = MapToPlayerInformationDTO(item)
		return &dto
	})
}

func mapToCardWithCountDTO(model *models.CardWithCount) *dtos.CardWithCountDTO {
	return &dtos.CardWithCountDTO{
		CardID:   model.Card.ID,
		CardName: model.Card.Name,
		Count:    model.Count,
	}
}

func mapToAttributeStrings(model []*dbModels.Attribute) []string {
	arrayCopy := make([]*dbModels.Attribute, len(model))
	copy(arrayCopy, model)
	slices.SortFunc(arrayCopy, func(attr0 *dbModels.Attribute, attr1 *dbModels.Attribute) int {
		return attr0.ID - attr1.ID
	})
	return lo.Map(
		arrayCopy,
		func(attr *dbModels.Attribute, _ int) string { return attr.Name },
	)
}

func MapToDeckDTO(model *models.Deck) *dtos.DeckDTO {
	return &dtos.DeckDTO{
		ID:         model.ID,
		Name:       model.Name,
		AvatarName: model.AvatarName,
		Cards:      lo.Map(model.Cards, func(card *models.CardWithCount, _ int) *dtos.CardWithCountDTO { return mapToCardWithCountDTO(card) }),
		Attributes: mapToAttributeStrings(model.Attributes),
	}
}

func MapToDeckDTOs(model []*models.Deck) []*dtos.DeckDTO {
	mapped := lo.Map(model, func(item *models.Deck, _ int) *dtos.DeckDTO { return MapToDeckDTO(item) })
	slices.SortFunc(mapped, func(attr0 *dtos.DeckDTO, attr1 *dtos.DeckDTO) int {
		return attr0.ID - attr1.ID
	})
	return mapped
}
