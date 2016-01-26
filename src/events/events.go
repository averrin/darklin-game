package events

import (
	"fmt"
	"time"
)

// EventType is a basic event type
type EventType int

// EventTypes enum
const (
	ALL EventType = iota
	INFO
	TICK
	SECOND
	MINUTE
	HOUR
	DAY
	MESSAGE
	COMMAND
	HEARTBEAT
	RESET //10
	CLOSE
	PAUSE
	ERROR
	LOGIN
	LOGGEDIN
	LOGINFAIL
	ROOMEXIT
	ROOMENTER
	ROOMCHANGED
	SYSTEMMESSAGE
	LIGHT
	CONNECTED
	DESCRIBE
)

// Event is atom of event stream
type Event struct {
	ID        string
	Timestamp time.Time
	Type      EventType
	Payload   interface{}
	Sender    string
	Delay     time.Duration
	Every     time.Duration
	Abort     bool
}

func (event Event) String() string {
	return fmt.Sprintf("{Sender: %v; Type: %v; Payload: %v, Abort: %v, ID: '%v', (Delay: %v, Every: %v)}",
		event.Sender, event.Type, event.Payload, event.Abort, event.ID, event.Delay, event.Every)
}

// NewEvent constructor
func NewEvent(eventType EventType, payload interface{}, sender string) *Event {
	event := new(Event)
	event.Timestamp = time.Now()
	event.Type = eventType
	event.Payload = payload
	event.Sender = sender
	event.Delay = 0
	event.Abort = false
	return event
}
