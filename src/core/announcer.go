package main

import "fmt"

// Announcer just someone who do something
type Announcer struct {
	Actor
}

// ConsumeEvent of couse
func (a Announcer) ConsumeEvent(event *Event) {
	a.Stream <- event
}

// NewAnnouncer because i, sucj in golang yet
func NewAnnouncer(gs *chan *Event) *Announcer {
	a := NewActor("Announcer", gs)
	actor := new(Announcer)
	actor.Actor = *a
	actor.Actor.ProcessEvent = actor.ProcessEvent
	return actor
}

// ProcessEvent - i need print something
func (a Announcer) ProcessEvent(event *Event) {
	switch event.Type {
	case SECOND:
		a.SendEvent("global", MESSAGE, "Every second, mister")
	case MINUTE:
		a.SendEvent("global", MESSAGE, fmt.Sprintf("Игровое время: %v", WORLD.Time.Date.Format("Mon Jan _2 15:04:05 2006")))
	}
}
