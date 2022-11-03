package chdist

import "encoding/json"

type JsonDistributor[T any] interface {
	Subscribe(f func(Item[T]))
	AsyncSubscribe() <-chan Item[T]
}

type jsonStringDistributor[T any] struct {
	Distributor Distributor[string]
}

func NewJsonStringDistributor[T any](distributor Distributor[string]) JsonDistributor[T] {
	return &jsonStringDistributor[T]{
		Distributor: distributor,
	}
}

func (j *jsonStringDistributor[T]) Subscribe(f func(Item[T])) {
	child := j.Distributor.AsyncSubscribe()
	go func() {
		for data := range child.Out() {
			instance := new(T)
			err := json.Unmarshal([]byte(data), &instance)
			f(Item[T]{*instance, err})
		}
	}()
}

func (j *jsonStringDistributor[T]) AsyncSubscribe() <-chan Item[T] {
	child := j.Distributor.AsyncSubscribe()
	ch := make(chan Item[T])
	go func() {
		for data := range child.Out() {
			instance := new(T)
			err := json.Unmarshal([]byte(data), &instance)
			ch <- Item[T]{*instance, err}
		}
	}()
	return ch
}

type jsonBytesDistributor[T any] struct {
	Distributor Distributor[[]byte]
}

func NewJsonBytesDistributor[T any](distributor Distributor[[]byte]) JsonDistributor[T] {
	return &jsonBytesDistributor[T]{
		Distributor: distributor,
	}
}

func (j *jsonBytesDistributor[T]) Subscribe(f func(Item[T])) {
	child := j.Distributor.AsyncSubscribe()
	go func() {
		for data := range child.Out() {
			instance := new(T)
			err := json.Unmarshal(data, &instance)
			f(Item[T]{*instance, err})
		}
	}()
}

func (j *jsonBytesDistributor[T]) AsyncSubscribe() <-chan Item[T] {
	child := j.Distributor.AsyncSubscribe()
	ch := make(chan Item[T])
	go func() {
		for data := range child.Out() {
			instance := new(T)
			err := json.Unmarshal(data, &instance)
			ch <- Item[T]{*instance, err}
		}
	}()
	return ch
}
