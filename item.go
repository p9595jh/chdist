package chdist

type Item[T any] struct {
	Value T
	Error error
}
