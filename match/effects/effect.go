package effects

import (
	"fmt"
	"strconv"

	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/models/enums"
)

type IEffect interface {
	GetType() enums.EffectType
	GetStringDescription() string
}

func NewEffect(dbEffect dbModels.CardEffect) (IEffect, error) {
	var effectType enums.EffectType = enums.EffectType(dbEffect.EffectID)
	switch effectType {
	case enums.EffectTypeCover, enums.EffectTypeShackled, enums.EffectTypeWounded:
		effect := NewEffectSimple(effectType, dbEffect.Name)
		return &effect, nil
	case enums.EffectTypeModifyPowerHealth:
		if dbEffect.Parameter0 == nil {
			return nil, fmt.Errorf("NewEffect: parameter0 is required for effect with a type '%d'", effectType)
		}
		p0, err := strconv.Atoi(*dbEffect.Parameter0)
		if err != nil {
			return nil, err
		}
		if dbEffect.Parameter1 == nil {
			return nil, fmt.Errorf("NewEffect: parameter1 is required for effect with a type '%d'", effectType)
		}
		p1, err := strconv.Atoi(*dbEffect.Parameter1)
		if err != nil {
			return nil, err
		}
		effect := NewEffectModifyPowerHealth(effectType, p0, p1)
		return &effect, nil
	default:
		return nil, fmt.Errorf("NewEffect: unknown effect type '%d'", effectType)
	}
}
