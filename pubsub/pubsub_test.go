package pubsub

import (
	"testing"
	"time"
)

func TestPublishDropMessage(t *testing.T) {
	p := New[[]byte]()

	p.Publish("non-existing", []byte("payload"))
}

func TestPubSubOK(t *testing.T) {
	t.Parallel()

	p := New[[]byte]()

	testPayload := "payload"
	testTopic := "topic1"

	testCh := p.Subscribe(testTopic)
	closeCh := make(chan struct{})

	go func() {
		msg := <-testCh
		if string(msg) != testPayload {
			t.Error("missmatch")
		}

		close(closeCh)
	}()

	time.Sleep(1 * time.Millisecond)
	p.Publish(testTopic, []byte(testPayload))

	<-closeCh
}

func TestPubSubMultiSubOK(t *testing.T) {
	t.Parallel()

	p := New[[]byte]()

	testPayload := "payload"
	testTopic := "topic1"

	testCh := p.Subscribe(testTopic)
	closeCh := make(chan struct{})

	test2Ch := p.Subscribe(testTopic)
	close2Ch := make(chan struct{})

	go func() {
		msg := <-testCh
		if string(msg) != testPayload {
			t.Error("missmatch")
		}

		close(closeCh)
	}()

	go func() {
		msg := <-test2Ch
		if string(msg) != testPayload {
			t.Error("missmatch")
		}

		close(close2Ch)
	}()

	time.Sleep(1 * time.Millisecond)
	p.Publish(testTopic, []byte(testPayload))

	<-closeCh
	<-close2Ch
}
