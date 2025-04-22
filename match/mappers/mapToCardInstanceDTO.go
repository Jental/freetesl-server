package mappers

import (
	"errors"

	dbEnums "github.com/jental/freetesl-server/db/enums"
	"github.com/jental/freetesl-server/match/dtos"
	"github.com/jental/freetesl-server/match/effects"
	"github.com/jental/freetesl-server/match/models"
	"github.com/samber/lo"
)

func MapToCardInstanceDTO(model models.CardInstance) (dtos.CardInstanceDTO, error) {
	switch cardInstanceCasted := model.(type) {
	case *models.CardInstanceCreature:
		return mapCreatureToCardInstanceDTO(cardInstanceCasted), nil
	case *models.CardInstanceItem:
		return mapItemToCardInstanceDTO(cardInstanceCasted), nil
	case *models.CardInstanceAction:
		return mapActionToCardInstanceDTO(cardInstanceCasted), nil
	case *models.CardInstanceSupport:
		return mapSupportToCardInstanceDTO(cardInstanceCasted), nil
	default:
		return dtos.CardInstanceDTO{}, errors.New("MapToCardInstanceDTO: Invalid cardInstance type")
	}
}

func mapToEffectInstanceDTO(model *effects.EffectInstance) dtos.EffectInstanceDTO {
	sourceCardInstanceID := model.SourceCardInstanceID.String()
	return dtos.EffectInstanceDTO{
		ID:                   byte(model.Effect.GetType()),
		Description:          model.Effect.GetStringDescription(),
		SourceCardInstanceID: &sourceCardInstanceID,
	}
}

func mapCreatureToCardInstanceDTO(model *models.CardInstanceCreature) dtos.CardInstanceDTO {
	return dtos.CardInstanceDTO{
		CardID:         model.Card.ID,
		CardInstanceID: model.CardInstanceID,
		Power:          model.GetComputedPower(),
		PowerMod:       model.GetPowerIncrease(),
		Health:         model.GetComputedHealth(),
		HealthMod:      model.GetHealthIncrease(),
		Cost:           model.Cost,
		Keywords:       lo.Map(model.GetAllKeywords(), func(kwd dbEnums.CardKeyword, _ int) int { return int(kwd) }),
		Effects: lo.Map(model.GetAllEffects(), func(eff *effects.EffectInstance, _ int) dtos.EffectInstanceDTO {
			return mapToEffectInstanceDTO(eff)
		}),
		// TODO:
		// - send unique effect types
	}
}

func mapItemToCardInstanceDTO(model *models.CardInstanceItem) dtos.CardInstanceDTO {
	return dtos.CardInstanceDTO{
		CardID:         model.Card.ID,
		CardInstanceID: model.CardInstanceID,
		Power:          0,
		Health:         0,
		Cost:           model.Cost,
		Keywords:       make([]int, 0),                    // there are no keywords valid for item, only definitions
		Effects:        make([]dtos.EffectInstanceDTO, 0), // there are no effect instances valid for item, only effect definitions
	}
}

func mapActionToCardInstanceDTO(model *models.CardInstanceAction) dtos.CardInstanceDTO {
	return dtos.CardInstanceDTO{
		CardID:         model.Card.ID,
		CardInstanceID: model.CardInstanceID,
		Power:          0,
		Health:         0,
		Cost:           model.Cost,
		Keywords:       lo.Map(model.Keywords, func(kwd dbEnums.CardKeyword, _ int) int { return int(kwd) }),
		Effects:        nil,
		// TODO:
		// - send unique effect types
		// - some effects (like silence) may overlap other effects - send only ones actual for FE
	}
}

func mapSupportToCardInstanceDTO(model *models.CardInstanceSupport) dtos.CardInstanceDTO {
	return dtos.CardInstanceDTO{
		CardID:         model.Card.ID,
		CardInstanceID: model.CardInstanceID,
		Power:          0,
		Health:         0,
		Cost:           model.Cost,
		Keywords:       lo.Map(model.Keywords, func(kwd dbEnums.CardKeyword, _ int) int { return int(kwd) }),
		Effects:        nil,
		// TODO:
		// - send unique effect types
		// - some effects (like silence) may overlap other effects - send only ones actual for FE
	}
}
