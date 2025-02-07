package common

type Maybe[T any] struct {
	Value    *T
	HasValue bool
}
