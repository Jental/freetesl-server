package match

import (
	"github.com/jental/freetesl-server/match/actions"
	"github.com/jental/freetesl-server/match/handlers"
	"github.com/jental/freetesl-server/match/interceptors"
	"github.com/jental/freetesl-server/match/match"
	"github.com/jental/freetesl-server/match/models"
	"github.com/jental/freetesl-server/match/senders"
	"github.com/jental/freetesl-server/models/enums"
)

func Main() {
	match.MatchMessageHandlerFn = handlers.ProcessMatchMessage
	match.BackendEventHandlerFn = senders.ProcessBackendEvent

	var guardInterceptor models.Interceptor = interceptors.GuardInterceptor{}
	interceptors.RegisterInterceptor(enums.InterceptorPointHitFaceBefore, &guardInterceptor)
	interceptors.RegisterInterceptor(enums.InterceptorPointHitCardBefore, &guardInterceptor)
	var coverInterceptor models.Interceptor = interceptors.CoverInterceptor{}
	interceptors.RegisterInterceptor(enums.InterceptorPointHitCardBefore, &coverInterceptor)

	interceptors.RegisterAllSpecialCardsInterceptors()

	actions.RegisterAllActions()
	actions.RegisterActionsForCards()
}
