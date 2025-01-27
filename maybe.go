package main

type Maybe[T any] struct {
	Value    T
	HasValue bool
}
