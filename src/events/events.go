package events

import "time"

// EventType is a basic event type
type EventType int

// EventTypes enum
const (
	ALL EventType = iota
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
}
