package models

type CardAction struct {
	CardID                 int
	ActionID               string
	InterceptorPointID     string
	ActionParametersValues *string // nullable
}
