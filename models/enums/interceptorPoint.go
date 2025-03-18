package enums

type InteceptorPoint string

const (
	InterceptorPointCardPlay                     InteceptorPoint = "operations.cardPlay"
	InterceptorPointMoveCardFromHandToLaneBefore InteceptorPoint = "operations.moveCardFromHandToLane:before"
	InterceptorPointMoveCardFromHandToLaneAfter  InteceptorPoint = "operations.moveCardFromHandToLane:after"
	InterceptorPointHitFaceBefore                InteceptorPoint = "operations.hitFace:before"
	InterceptorPointHitFaceAfter                 InteceptorPoint = "operations.hitFace:after"
)
