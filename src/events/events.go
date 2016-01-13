package events

import "time"

// EventType is a basic event type
type EventType int

// EventTypes enum
const (
	TICK EventType = iota
	MESSAGE
)

// Event is atom of event stream
type Event struct {
	Timestamp time.Time
	Type      EventType
	Payload   interface{}
}
