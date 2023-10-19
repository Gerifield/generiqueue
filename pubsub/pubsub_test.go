package pubsub

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, testPayload, string(msg))

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
		assert.Equal(t, testPayload, string(msg))

		close(closeCh)
	}()

	go func() {
		msg := <-test2Ch
		assert.Equal(t, testPayload, string(msg))

		close(close2Ch)
	}()

	time.Sleep(1 * time.Millisecond)
	p.Publish(testTopic, []byte(testPayload))

	<-closeCh
	<-close2Ch
}

func TestFinishCleanup(t *testing.T) {
	t.Parallel()

	p := New[[]byte]()

	testTopic := "topic1"
	testTopic2 := "topic2"

	_ = p.Subscribe(testTopic)
	_ = p.Subscribe(testTopic)
	_ = p.Subscribe(testTopic2)

	assert.Equal(t, 2, len(p.topics))
	assert.Equal(t, 2, len(p.topics[testTopic]))

	p.Finish(testTopic)
	assert.Equal(t, 1, len(p.topics))
	_, ok := p.topics[testTopic]
	assert.False(t, ok)
}

func TestPubSubBufferedOK(t *testing.T) {
	t.Parallel()

	p := New[[]byte]()

	testPayload := "payload"
	testTopic := "topic1"

	testCh := p.SubscribeBuffered(testTopic, 4)
	closeCh := make(chan struct{})

	go func() {
		cnt := 0
		for msg := range testCh {
			assert.Equal(t, testPayload, string(msg))
			cnt++
		}

		assert.Equal(t, 4, cnt)
		close(closeCh)
	}()

	time.Sleep(1 * time.Millisecond)
	// Not the best way to test, but without the buffer some of these would be dropped
	p.Publish(testTopic, []byte(testPayload))
	p.Publish(testTopic, []byte(testPayload))
	p.Publish(testTopic, []byte(testPayload))
	p.Publish(testTopic, []byte(testPayload))
	p.Finish(testTopic)

	<-closeCh
}

func TestPubSubCombined(t *testing.T) {
	t.Parallel()

	p := New[[]byte]()

	testPayload := "payload"
	testTopic := "topic1"

	var wg sync.WaitGroup
	testCh1 := p.SubscribeBuffered(testTopic, 4) // Buffered
	testCh2 := p.Subscribe(testTopic)            // Non buffered

	go func() {
		wg.Add(1)
		cnt := 0
		for msg := range testCh1 {
			assert.Equal(t, testPayload, string(msg))
			cnt++
		}

		assert.Equal(t, 4, cnt)
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		cnt := 0
		for msg := range testCh2 {
			assert.Equal(t, testPayload, string(msg))
			cnt++
		}

		// Without the buffer some of the events would be dropped
		assert.Less(t, cnt, 4)
		wg.Done()
	}()

	time.Sleep(1 * time.Millisecond)
	p.Publish(testTopic, []byte(testPayload))
	p.Publish(testTopic, []byte(testPayload))
	p.Publish(testTopic, []byte(testPayload))
	p.Publish(testTopic, []byte(testPayload))
	p.Finish(testTopic)

	wg.Wait()
}
