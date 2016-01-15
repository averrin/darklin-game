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
	Subscriptions []Subscription
	Streams       map[string]chan events.Event
	Name          string
	ID            string
}

// NewActor construct new Actor
func NewActor(name string, gs chan events.Event) *Actor {
	actor := new(Actor)
	actor.Streams = make(map[string]chan events.Event)
	actor.Streams["global"] = gs
	actor.Stream = make(chan events.Event)
	actor.Name = name
	return actor
}

// SendEvent with type and payload
func (a Actor) SendEvent(reciever string, eventType events.EventType, payload interface{}) {
	event := events.Event{
		time.Now(),
		eventType,
		payload,
		a.Name,
	}
	stream := a.Streams[reciever]
	stream <- event
}

// ForwardEvent to new reciever
func (a Actor) ForwardEvent(reciever string, event events.Event) {
	a.Streams[reciever] <- event
}

// Subscribe on events
func (a *Actor) Subscribe(eventType events.EventType, subscriber EventSubscriber) {
	a.Subscriptions = append(a.Subscriptions, Subscription{eventType, subscriber})
}

// NotifySubscribers wgen u have event
func (a Actor) NotifySubscribers(event events.Event) {
	for _, s := range a.Subscriptions {
		if event.Type == s.Type || s.Type == events.ALL {
			go s.Subscriber.ConsumeEvent(event)
		}
	}
}

// AddStream to Streams
func (a *Actor) AddStream(subscriber Actor) {
	a.Streams[subscriber.Name] = subscriber.Stream
}
