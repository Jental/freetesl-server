package mappers

import (
	"errors"

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
	var sourceCardInstanceIDStr *string = nil
	if model.SourceCardInstanceID != nil {
		str := model.SourceCardInstanceID.String()
		sourceCardInstanceIDStr = &str
	}
	return dtos.EffectInstanceDTO{
		ID:                   byte(model.Effect.GetType()),
		Description:          model.Effect.GetStringDescription(),
		SourceCardInstanceID: sourceCardInstanceIDStr,
	}
}

func mapToKeywordInstanceDTO(model *models.KeywordInstance) dtos.KeywordInstanceDTO {
	sourceCardInstanceIDStr := model.SourceCardInstanceID.String()
	return dtos.KeywordInstanceDTO{
		ID:                   byte(model.Keyword),
		SourceCardInstanceID: &sourceCardInstanceIDStr,
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
		Keywords: lo.Map(model.GetAllKeywords(), func(kwd *models.KeywordInstance, _ int) dtos.KeywordInstanceDTO {
			return mapToKeywordInstanceDTO(kwd)
		}),
		Effects: lo.Map(model.GetAllEffects(), func(eff *effects.EffectInstance, _ int) dtos.EffectInstanceDTO {
			return mapToEffectInstanceDTO(eff)
		}),
		// TODO: maybe send unique effect types
	}
}

func mapItemToCardInstanceDTO(model *models.CardInstanceItem) dtos.CardInstanceDTO {
	return dtos.CardInstanceDTO{
		CardID:         model.Card.ID,
		CardInstanceID: model.CardInstanceID,
		Power:          0,
		Health:         0,
		Cost:           model.Cost,
		Keywords:       make([]dtos.KeywordInstanceDTO, 0), // there are no keywords valid for item, only definitions
		Effects:        make([]dtos.EffectInstanceDTO, 0),  // there are no effect instances valid for an item, only effect definitions
	}
}

func mapActionToCardInstanceDTO(model *models.CardInstanceAction) dtos.CardInstanceDTO {
	return dtos.CardInstanceDTO{
		CardID:         model.Card.ID,
		CardInstanceID: model.CardInstanceID,
		Power:          0,
		Health:         0,
		Cost:           model.Cost,
		Keywords: lo.Map(model.KeywordInstances, func(kwd *models.KeywordInstance, _ int) dtos.KeywordInstanceDTO {
			return mapToKeywordInstanceDTO(kwd)
		}),
		Effects: make([]dtos.EffectInstanceDTO, 0), // there are no effect instances valid for an action
	}
}

func mapSupportToCardInstanceDTO(model *models.CardInstanceSupport) dtos.CardInstanceDTO {
	return dtos.CardInstanceDTO{
		CardID:         model.Card.ID,
		CardInstanceID: model.CardInstanceID,
		Power:          0,
		Health:         0,
		Cost:           model.Cost,
		Keywords: lo.Map(model.KeywordInstances, func(kwd *models.KeywordInstance, _ int) dtos.KeywordInstanceDTO {
			return mapToKeywordInstanceDTO(kwd)
		}),
		Effects: make([]dtos.EffectInstanceDTO, 0), // TODO: supports can have effects, e.g. Last Gasp
	}
}
