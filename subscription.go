package chdist

type Subscription[T any] interface {
	Close()
	IsAlive() bool
	Subject() Distributor[T]
}

type AsyncSubscription[T any] interface {
	Subscription[T]
	Out() <-chan T
}

type subscription[T any] struct {
	channel chan T
	isAlive bool
	id      uint64
	parent  *distributor[T]
}

func (s *subscription[T]) Close() {
	s.isAlive = false
	delete(s.parent.children, s.id)
	close(s.channel)
}

func (s *subscription[T]) IsAlive() bool {
	return s.isAlive
}

func (s *subscription[T]) Subject() Distributor[T] {
	return s.parent
}

func (s *subscription[T]) Out() <-chan T {
	return s.channel
}
