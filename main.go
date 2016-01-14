package main

import (
	"actor"
	"events"
	"global"
	// "time"
	"time_stream"
	// "fmt"
)

// TestActor just someone who do something
type TestActor struct {
	actor.Actor
}

// ConsumeEvent of couse
func (a TestActor) ConsumeEvent(event events.Event) {
	a.Stream <- event
}

// NewTestActor because i, sucj in golang yet
func NewTestActor(gs chan events.Event) *TestActor {
	actor := new(TestActor)
	actor.GlobalStream = gs
	actor.Stream = make(chan events.Event)
	return actor
}

// Live - i need print something
func (a TestActor) Live() {
	for {
		event := <- a.Stream
		for _, s := range a.Subscriptions {
			if event.Type == s.Type {
				go s.Subscriber.ConsumeEvent(event)
			}
		}
		if event.Timestamp.Second() % 2 == 0 {
			a.SendEvent(events.MESSAGE, "even Tick")
		}
	}
}

func main() {
	stream := make(chan events.Event)
	ts := time_stream.TimeStream{}
	ts.GlobalStream = stream
	go ts.Live()

	testActor := NewTestActor(stream)
	go testActor.Live()

	gs := global.Stream{}
	gs.Stream = stream
	gs.Subscribe(events.TICK, testActor)
	gs.Live()
}
