package models

import (
	"github.com/jental/freetesl-server/models/enums"
)

type InterceptorContext struct {
	PlayerState   *PlayerMatchState
	OpponentState *PlayerMatchState
	LaneID        *enums.Lane
}

type Interceptor interface {
	GetInterceptorPoint() enums.InteceptorPoint
	Execute(context *InterceptorContext) error
}
