package main

import (
	"fmt"
	"github.com/gerifield/mini-pubsub/pubsub"
	"time"
)

func main() {
	p := pubsub.New()

	msgCh := p.Subscribe("topic1")

	stopCh := make(chan struct{})
	go func() {
		for v := range msgCh {
			fmt.Println("Message:", string(v))
		}

		close(stopCh)
	}()
	time.Sleep(20 * time.Millisecond)

	p.Publish("topic1", []byte("hello world"))
	time.Sleep(20 * time.Millisecond)
	p.Publish("topic1", []byte("hello world2"))
	p.Finish("topic1")

	<-stopCh
}
