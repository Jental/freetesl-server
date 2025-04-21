package effects

import "github.com/jental/freetesl-server/models/enums"

type EffectSimple struct {
	effectType        enums.EffectType
	stringDescription string
}

func (effect *EffectSimple) GetType() enums.EffectType {
	return effect.effectType
}

func (effect *EffectSimple) GetStringDescription() string {
	return effect.stringDescription
}

func NewEffectSimple(effectType enums.EffectType, stringDescription string) EffectSimple {
	return EffectSimple{
		effectType:        effectType,
		stringDescription: stringDescription,
	}
}
