package interceptors

import (
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
)

var interceptors map[enums.InteceptorPoint][]*models.Interceptor = make(map[enums.InteceptorPoint][]*models.Interceptor)

func ExecuteInterceptors(point enums.InteceptorPoint, context *models.InterceptorContext) error {
	pointIcs, exists := interceptors[point]
	if !exists {
		return nil
	}

	for _, ic := range pointIcs {
		err := (*ic).Execute(context)
		if err != nil {
			return err
		}
	}

	return nil
}

func RegisterInterceptor(point enums.InteceptorPoint, interceptor *models.Interceptor) {
	ics, exists := interceptors[point]
	if !exists {
		ics = make([]*models.Interceptor, 0)
	}
	interceptors[point] = append(ics, interceptor)
}

func RegisterAllSpecialCardsInterceptors() {
	var brumaProfiteerInterceptor models.Interceptor = BrumaProfiteerInterceptor{}
	RegisterInterceptor(enums.InterceptorPointMoveCardFromHandToLaneAfter, &brumaProfiteerInterceptor)
}
