package events

import "time"

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
	MESSAGE
	COMMAND
)

// Event is atom of event stream
type Event struct {
	Timestamp time.Time
	Type      EventType
	Payload   interface{}
	Sender    string
}

// NewEvent constructor
func NewEvent(eventType EventType, payload interface{}, sender string) Event {
	event := new(Event)
	event.Timestamp = time.Now()
	event.Type = eventType
	event.Payload = payload
	event.Sender = sender
	return *event
}
