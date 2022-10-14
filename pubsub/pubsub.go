package pubsub

import (
	"sync"
)

// PubSub .
type PubSub[T any] struct {
	incomingLock sync.Mutex
	topics       map[string][]chan T
}

// New .
func New[T any]() *PubSub[T] {
	return &PubSub[T]{
		topics: make(map[string][]chan T),
	}
}

// Publish .
func (ps *PubSub[T]) Publish(topic string, payload T) {
	ps.incomingLock.Lock()
	defer ps.incomingLock.Unlock()

	if targetChs, ok := ps.topics[topic]; ok {
		for _, ch := range targetChs {
			select {
			case ch <- payload:
			default:
			}
		}
	}
}

// Finish .
func (ps *PubSub[T]) Finish(topic string) {
	ps.incomingLock.Lock()
	defer ps.incomingLock.Unlock()

	if targetChs, ok := ps.topics[topic]; ok {
		for _, ch := range targetChs {
			close(ch)
		}
	}

	delete(ps.topics, topic)
}

// Subscribe .
func (ps *PubSub[T]) Subscribe(topic string) <-chan T {
	ps.incomingLock.Lock()
	defer ps.incomingLock.Unlock()

	retCh := make(chan T)

	_, ok := ps.topics[topic]
	if ok {
		ps.topics[topic] = append(ps.topics[topic], retCh)

		return retCh
	}

	ps.topics[topic] = []chan T{retCh}

	return retCh
}
