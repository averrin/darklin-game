package main

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
	MESSAGE
	COMMAND
	HEARTBEAT //8
	RESET
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
)

// Event is atom of event stream
type Event struct {
	Timestamp time.Time
	Type      EventType
	Payload   interface{}
	Sender    string
}

func (event Event) String() string {
	return fmt.Sprintf("{Sender: %v; Type: %v; Payload: %v}", event.Sender, event.Type, event.Payload)
}

// NewEvent constructor
func NewEvent(eventType EventType, payload interface{}, sender string) *Event {
	event := new(Event)
	event.Timestamp = time.Now()
	event.Type = eventType
	event.Payload = payload
	event.Sender = sender
	return event
}
