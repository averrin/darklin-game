package main

import (
	"actor"
	"events"
	"global"
	"time"
)

type TestActor struct {
	actor.Actor
}

func (actor TestActor) ConsumeEvent(event events.Event) {
	actor.Stream <- event
}

func (a TestActor) Live() {
	for {
		event := <-a.Stream
		for _, s := range a.Subscriptions {
			if event.Type == s.Type {
				go s.Subscriber.ConsumeEvent(event)
			}
		}
		if event.Timestamp.Second()%2 == 0 {
			a.SendEvent(events.MESSAGE, "even Tick")
		}
	}
}

func main() {
	stream := make(chan events.Event)
	go timeStream(stream)
	testActor := TestActor{actor.NewActor(stream)}
	testActor.Stream = stream
	go testActor.Live()

	var subscriptions []actor.Subscription
	gs := global.GlobalStream{}
	gs.Actor.Stream = stream
	gs.Subscriptions = append(subscriptions, actor.Subscription{events.TICK, testActor})
	gs.Live()
}

func timeStream(stream chan events.Event) {
	ticker := time.NewTicker(time.Millisecond * 500)
	for t := range ticker.C {
		stream <- events.Event{t, events.TICK, nil}
	}
}
