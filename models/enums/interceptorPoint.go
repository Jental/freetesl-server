package enums

type InteceptorPoint string

const (
	InterceptorPointHitFaceBefore InteceptorPoint = "operations.hitFace:before"
	InterceptorPointHitFaceAfter  InteceptorPoint = "operations.hitFace:after"
)
