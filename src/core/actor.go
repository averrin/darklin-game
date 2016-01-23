package main

import "log"

// "fmt"

// Interface - Anybody who can live
type Interface interface {
	Live()
	ProcessEvent(event *Event)
}

// EventSubscriber - can subscribe
type EventSubscriber interface {
	ConsumeEvent(*Event)
}

// EventPublisher - can send
type EventPublisher interface {
	SendEvent(EventType, interface{})
}

// Subscription on events
type Subscription struct {
	Type       EventType
	Subscriber EventSubscriber
}

// Actor - basic event-driven class
type Actor struct {
	Stream         chan *Event
	Subscriptions  []Subscription
	Streams        map[string]*chan *Event
	Name           string
	ID             string
	Storage        *Storage
	ProcessEvent   func(event *Event)
	ProcessCommand func(event *Event)
}

// NewActor construct new Actor
func NewActor(name string, gs *chan *Event) *Actor {
	actor := new(Actor)
	actor.Streams = make(map[string]*chan *Event)
	actor.Streams["global"] = gs
	actor.Stream = make(chan *Event)
	actor.Name = name
	actor.Storage = NewStorage()
	actor.ProcessEvent = actor.ProcessEventAbstract
	actor.ProcessCommand = actor.ProcessCommandAbstract
	return actor
}

// SendEvent with type and payload
func (a Actor) SendEvent(reciever string, eventType EventType, payload interface{}) {
	event := NewEvent(eventType, payload, a.Name)
	stream := a.Streams[reciever]
	*stream <- event
}

// SendEventWithSender - fake sender
func (a Actor) SendEventWithSender(reciever string, eventType EventType, payload interface{}, sender string) {
	event := NewEvent(eventType, payload, sender)
	stream := a.Streams[reciever]
	*stream <- event
}

// Broadcast - send all
func (a Actor) Broadcast(eventType EventType, payload interface{}, sender string) {
	event := NewEvent(eventType, payload, sender)
	defer func() { recover() }()
	// yellow := color.New(color.FgYellow).SprintFunc()
	// if event.Type != HEARTBEAT {
	// 	log.Println(yellow("Broadcast event"), event)
	// }
	for r := range a.Streams {
		if r == "global" || r == sender || r == "time" {
			continue
		}
		*a.Streams[r] <- event
	}
}

// BroadcastRoom - send all
func (a *Actor) BroadcastRoom(eventType EventType, payload interface{}, sender string, room *Area) {
	event := NewEvent(eventType, payload, sender)
	defer func() { recover() }()
	for p := range room.Players {
		if p.Name == sender {
			continue
		}
		p.Stream <- event
	}
}

// ForwardEvent to new reciever
func (a Actor) ForwardEvent(reciever string, event *Event) {
	// defer func() { recover() }()
	log.Println("event before forwarded", reciever, *a.Streams[reciever])
	*a.Streams[reciever] <- event
	log.Println("event forwarded")
}

// Subscribe on events
func (a *Actor) Subscribe(eventType EventType, subscriber EventSubscriber) {
	a.Subscriptions = append(a.Subscriptions, Subscription{eventType, subscriber})
}

// NotifySubscribers wgen u have event
func (a Actor) NotifySubscribers(event *Event) {
	for _, s := range a.Subscriptions {
		if event.Type == s.Type || s.Type == ALL {
			go s.Subscriber.ConsumeEvent(event)
		}
	}
}

// AddStream to Streams
func (a *Actor) AddStream(subscriber Actor) {
	a.Streams[subscriber.Name] = &subscriber.Stream
}

// Live method for dispatch events
func (a *Actor) Live() {
	s := a.Storage.Session.Copy()
	defer s.Close()
	a.Storage.DB = s.DB("darklin")
	for {
		event := <-a.Stream
		exp_events_processed.Add(1)
		// log.Println(a.Name, event)
		a.NotifySubscribers(event)
		a.ProcessEvent(event)
	}
	// log.Println(a.Formatter.Red("Live stopped"))
}

//ProcessEventAbstract - dummy processor
func (a *Actor) ProcessEventAbstract(event *Event) {
	log.Println("Abstract", a.Name, event)
}

func (a *Actor) ProcessCommandAbstract(event *Event) {}
