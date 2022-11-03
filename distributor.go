package chdist

import (
	"fmt"
	"strconv"
)

type Distributor[T any] interface {
	In() chan<- T
	Len() int
	Close()
	Subscriptions() []Subscription[T]
	Subscribe(f func(value T)) Subscription[T]
	AsyncSubscribe() AsyncSubscription[T]
}

type distributor[T any] struct {
	channel  chan T
	children map[uint64]*subscription[T]
}

func NewDistributor[T any](channel chan T) Distributor[T] {
	dist := &distributor[T]{
		channel:  channel,
		children: make(map[uint64]*subscription[T]),
	}
	go func() {
		for data := range channel {
			for _, child := range dist.children {
				child.channel <- data
			}
		}
	}()
	return dist
}

func (d *distributor[T]) In() chan<- T {
	return d.channel
}

func (d *distributor[T]) Len() int {
	return len(d.children)
}

func (d *distributor[T]) Subscriptions() []Subscription[T] {
	subscriptions := make([]Subscription[T], len(d.children))
	i := 0
	for _, child := range d.children {
		subscriptions[i] = child
		i++
	}
	return subscriptions
}

func (d *distributor[T]) insert() uint64 {
	child := &subscription[T]{
		channel: make(chan T),
		isAlive: true,
		parent:  d,
	}

	id, _ := strconv.ParseUint(fmt.Sprintf("%p", &child)[2:], 16, 64)
	child.id = id
	d.children[id] = child
	return id
}

func (d *distributor[T]) Subscribe(f func(value T)) Subscription[T] {
	child := d.children[d.insert()]
	go func() {
		for data := range child.channel {
			f(data)
		}
	}()
	return child
}

func (d *distributor[T]) AsyncSubscribe() AsyncSubscription[T] {
	child := d.children[d.insert()]
	return child
}

func (d *distributor[T]) Close() {
	for _, child := range d.children {
		child.Close()
	}
	close(d.channel)
}
