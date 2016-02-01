package actor

import (
	"core"
	"events"
	"expvar"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var (
	exp_events_processed = expvar.NewInt("events_processed")
)

// "fmt"

// Subscription on events
type Subscription struct {
	Type       events.EventType
	Subscriber *Actor
}

//CharState - Basic state
type CharState struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string

	Room string
	HP   int

	New bool
}

//AreaState - db saved state
type AreaState struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string

	Light bool
	Items []string

	New bool
}

// Actor - basic event-driven class
type Actor struct {
	Stream        chan *events.Event
	Subscriptions []Subscription
	Streams       map[string]*chan *events.Event
	Name          string
	ID            string
	Storage       *core.Storage
	World         WorldInterface

	PendingEvents map[string]*events.Event

	Handlers        map[events.EventType]func(*events.Event) bool
	CommandHandlers map[string]func(string) bool
	ProcessEvent    func(event *events.Event)
	ProcessCommand  func(event *events.Event)
}

//String func for plain actor
// func (a *Actor) String() string {
// 	return fmt.Sprintf("{Name: %s}", a.Name)
// }

// NewActor construct new Actor
func NewActor(name string, gs StreamInterface) Actor {
	actor := new(Actor)
	actor.World = gs.GetWorld()
	actor.Streams = make(map[string]*chan *events.Event)
	actor.PendingEvents = make(map[string]*events.Event)
	actor.Streams["global"] = gs.GetStream()
	actor.Stream = make(chan *events.Event)
	actor.Name = name
	actor.Storage = core.NewStorage()
	actor.Handlers = make(map[events.EventType]func(*events.Event) bool)
	actor.CommandHandlers = make(map[string]func(string) bool)
	actor.ProcessEvent = actor.ProcessEventAbstract
	actor.ProcessCommand = actor.ProcessCommandAbstract
	return *actor
}

// SendEvent with type and payload
func (a *Actor) SendEvent(reciever string, eventType events.EventType, payload interface{}) {
	event := events.NewEvent(eventType, payload, a.Name)
	stream := a.Streams[reciever]
	*stream <- event
}

// SendEventWithSender - fake sender
func (a *Actor) SendEventWithSender(reciever string, eventType events.EventType, payload interface{}, sender string) {
	// log.Println(a.Streams)
	event := events.NewEvent(eventType, payload, sender)
	stream := a.Streams[reciever]
	*stream <- event
}

// Broadcast - send all
func (a *Actor) Broadcast(eventType events.EventType, payload interface{}, sender string) {
	event := events.NewEvent(eventType, payload, sender)
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

// ForwardEvent to new reciever
func (a *Actor) ForwardEvent(reciever string, event *events.Event) {
	// defer func() { recover() }()
	// log.Println("event before forwarded", reciever, *a.Streams[reciever])
	*a.Streams[reciever] <- event
	// log.Println("event forwarded")
}

// Subscribe on events
func (a *Actor) Subscribe(eventType events.EventType, subscriber *Actor) {
	a.Subscriptions = append(a.Subscriptions, Subscription{eventType, subscriber})
}

// NotifySubscribers wgen u have event
func (a *Actor) NotifySubscribers(event *events.Event) {
	for _, s := range a.Subscriptions {
		if event.Type == s.Type || s.Type == events.ALL {
			s.Subscriber.Stream <- event
		}
	}
}

// AddStream to Streams
func (a *Actor) AddStream(subscriber Actor) {
	a.Streams[subscriber.Name] = &subscriber.Stream
}

// Live method for dispatch events
func (a *Actor) Live() {
	// if a.Name == "Hall" {
	// log.Println(a)
	// }
	s := a.Storage.Session.Copy()
	defer s.Close()
	a.Storage.DB = s.DB("darklin")
	for {
		event := <-a.Stream
		exp_events_processed.Add(1)
		if event.Abort {
			continue
		}
		if event.Delay != 0 {
			if event.ID != "" {
				a.PendingEvents[event.ID] = event
			}
			go func() {
				a.Sleep(event.Delay)
				event.Delay = 0
				if event.ID != "" {
					delete(a.PendingEvents, event.ID)
				}
				a.Stream <- event
			}()
			continue
		}
		if event.Every != 0 {
			go func() {
				for !event.Abort {
					a.Sleep(event.Every)
					a.NotifySubscribers(event)
					a.ProcessEvent(event)
				}
			}()
			continue
		}
		// log.Println(a.Name, event)
		a.NotifySubscribers(event)
		a.ProcessEvent(event)
	}
	// log.Println(a.Formatter.Red("Live stopped"))
}

//ProcessEventAbstract - dummy processor
func (a *Actor) ProcessEventAbstract(event *events.Event) {
	log.Println("Abstract", a.Name, event)
}

//ProcessCommandAbstract -
func (a *Actor) ProcessCommandAbstract(event *events.Event) {}

//Sleep -
func (a *Actor) Sleep(duration time.Duration) {
	(*a.World.GetTime()).Sleep(duration)
}

//GetName -
func (a *Actor) GetName() string {
	return a.Name
}

//GetStream -
func (a *Actor) GetStream() *chan *events.Event {
	return &a.Stream
}

//GetPendingEvent -
func (a *Actor) GetPendingEvent(name string) (*events.Event, bool) {
	ev, ok := a.PendingEvents[name]
	return ev, ok
}

//SetWorld -
func (a *Actor) SetWorld(w WorldInterface) {
	a.World = w
}

//SetStream -
func (a *Actor) SetStream(name string, s *chan *events.Event) {
	a.Streams[name] = s
}

//GetWorld -
func (a *Actor) GetWorld() WorldInterface {
	return a.World
}

//Index -
func Index(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}
