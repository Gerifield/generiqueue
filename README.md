# generiqueue
Mini experimental in-memory pubsub with generics

I just wanted to experiment a bit with the generics to at least try them and this seemd as a good little example.

## Usage

You can create a simple pubsub instance with a given type it'll transfer, for example for `[]byte`:

```go
p := pubsub.New[[]byte]()
```

Then you can just publish the message.
```go
p.Publish("topic1", []byte("hello world"))
```

Good to know, this will always work, it does not block and will drop the message if there's no subscriber!

The subscribe call will return a channel where the messages will come, you can use a loop to read them:

```go
msgCh := p.Subscribe("topic1")
for m := range msgCh {
	fmt.Println(m)
}
```

You can have multiple subscribers on the same topic, they'll receive a copy of the message each.

There's an additional method to signal the subscribers there won't be any new message in the given channel:

```go
p.Finish("topic1")
```

This will close the channel and with that "stop" the range loops on the subscribers.
