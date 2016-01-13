package actor

import (
	"events"
	"time"
)

type EventSubscriber interface {
	ConsumeEvent(events.Event)
	Live()
}
type EventPublisher interface {
	SendEvent(events.EventType, interface{})
}

type Subscription struct {
	Type       events.EventType
	Subscriber EventSubscriber
}

type Actor struct {
	Stream        chan events.Event
	GlobalStream  chan events.Event
	Subscriptions []Subscription
}

// Actor constructor
func NewActor(gs chan events.Event) Actor {
	actor := new(Actor)
	actor.GlobalStream = gs
	actor.Stream = make(chan events.Event)
	return *actor
}

func (a Actor) SendEvent(eventType events.EventType, payload interface{}) {
	event := events.Event{
		time.Now(),
		eventType,
		payload,
	}
	a.GlobalStream <- event
}
