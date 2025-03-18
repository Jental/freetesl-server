package models

type Interceptor interface {
	Execute(context *InterceptorContext) error
}
