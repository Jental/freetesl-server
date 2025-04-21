package effects

import (
	"fmt"

	"github.com/jental/freetesl-server/models/enums"
)

type EffectModifyPowerHealth struct {
	effectType     enums.EffectType
	PowerIncrease  int
	HealthIncrease int
}

func (effect *EffectModifyPowerHealth) GetType() enums.EffectType {
	return effect.effectType
}

func buildModifyString(val int) string {
	if val == 0 {
		return "0"
	} else if val > 0 {
		return fmt.Sprintf("+%d", val)
	} else {
		return fmt.Sprintf("-%d", val)
	}
}

func (effect *EffectModifyPowerHealth) GetStringDescription() string {
	powerIncreaseStr := buildModifyString(effect.PowerIncrease)
	healthIncreaseStr := buildModifyString(effect.HealthIncrease)
	return fmt.Sprintf("%s/%s", powerIncreaseStr, healthIncreaseStr)
}

func NewEffectModifyPowerHealth(effectType enums.EffectType, powerIncrease int, healthIncrease int) EffectModifyPowerHealth {
	return EffectModifyPowerHealth{
		effectType:     effectType,
		PowerIncrease:  powerIncrease,
		HealthIncrease: healthIncrease,
	}
}
