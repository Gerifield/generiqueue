package pubsub

import (
	"sync"
)

// PubSub .
type PubSub struct {
	incomingLock sync.Mutex
	topics       map[string][]chan []byte
}

// New .
func New() *PubSub {
	return &PubSub{
		topics: make(map[string][]chan []byte),
	}
}

// Publish .
func (ps *PubSub) Publish(topic string, payload []byte) {
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
func (ps *PubSub) Finish(topic string) {
	ps.incomingLock.Lock()
	defer ps.incomingLock.Unlock()

	if targetChs, ok := ps.topics[topic]; ok {
		for _, ch := range targetChs {
			close(ch)
		}
	}
}

// Subscribe .
func (ps *PubSub) Subscribe(topic string) <-chan []byte {
	ps.incomingLock.Lock()
	defer ps.incomingLock.Unlock()

	retCh := make(chan []byte)

	_, ok := ps.topics[topic]
	if ok {
		ps.topics[topic] = append(ps.topics[topic], retCh)

		return retCh
	}

	ps.topics[topic] = []chan []byte{retCh}

	return retCh
}
