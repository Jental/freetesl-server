package enums

type InteceptorPoint string

const (
	InterceptorPointCardPlay             InteceptorPoint = "operations.cardPlay"
	InterceptorPointMoveCardToLaneBefore InteceptorPoint = "operations.moveCardToLane:before"
	InterceptorPointMoveCardToLaneAfter  InteceptorPoint = "operations.moveCardToLane:after"
	InterceptorPointHitFaceBefore        InteceptorPoint = "operations.hitFace:before"
	InterceptorPointHitFaceAfter         InteceptorPoint = "operations.hitFace:after"
	InterceptorPointHitCardBefore        InteceptorPoint = "operations.hitCard:before"
	InterceptorPointHitCardAfter         InteceptorPoint = "operations.hitCard:after"
	InterceptorPointHealthReduceBefore   InteceptorPoint = "operations.reducePlayerHealth:before"
	InterceptorPointHealthReduceAfter    InteceptorPoint = "operations.reducePlayerHealth:after"
	InterceptorPointRuneBreakBefore      InteceptorPoint = "operations.reducePlayerHealth.runeBreak:before"
	InterceptorPointRuneBreakAfter       InteceptorPoint = "operations.reducePlayerHealth.runeBreak:after"
)
