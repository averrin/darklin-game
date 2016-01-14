package actor

import (
	"events"
	"time"
	// "fmt"
)

// Interface - Anybody who can live
type Interface interface {
	Live()
}

// EventSubscriber - can subscribe
type EventSubscriber interface {
	ConsumeEvent(events.Event)
}

// EventPublisher - can send
type EventPublisher interface {
	SendEvent(events.EventType, interface{})
}

// Subscription on events
type Subscription struct {
	Type       events.EventType
	Subscriber EventSubscriber
}

// Actor - basic event-driven class
type Actor struct {
	Stream        chan events.Event
	GlobalStream  chan events.Event
	Subscriptions []Subscription
}

// NewActor construct new Actor
func NewActor(gs chan events.Event) *Actor {
	actor := new(Actor)
	actor.GlobalStream = gs
	actor.Stream = make(chan events.Event)
	return actor
}

// SendEvent with type and payload
func (a Actor) SendEvent(eventType events.EventType, payload interface{}) {
	event := events.Event{
		time.Now(),
		eventType,
		payload,
	}
	a.GlobalStream <- event
}

// Subscribe on events
func (a *Actor) Subscribe(eventType events.EventType, subscriber EventSubscriber) {
	a.Subscriptions = append(a.Subscriptions, Subscription{eventType, subscriber})
}
