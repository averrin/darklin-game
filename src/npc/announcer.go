package npc

import "fmt"
import "actor"
import "events"
import "world"

// Announcer just someone who do something
type Announcer struct {
	actor.Actor
}

// ConsumeEvent of couse
func (a Announcer) ConsumeEvent(event *events.Event) {
	a.Stream <- event
}

// NewAnnouncer because i, sucj in golang yet
func NewAnnouncer(gs *chan *events.Event) *Announcer {
	a := actor.NewActor("Announcer", gs)
	announcer := new(Announcer)
	announcer.Actor = *a
	announcer.Actor.ProcessEvent = announcer.ProcessEvent
	return announcer
}

// ProcessEvent - i need print something
func (a Announcer) ProcessEvent(event *events.Event) {
	switch event.Type {
	case SECOND:
		a.SendEvent("global", events.MESSAGE, "Every second, mister")
	case MINUTE:
		a.SendEvent("global", events.MESSAGE, fmt.Sprintf("Игровое время: %v", world.WORLD.Time.Date.Format("Mon Jan _2 15:04:05 2006")))
	}
}
