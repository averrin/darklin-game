package npc

import "fmt"
import actor "../actor"
import events "../events"

// Announcer just someone who do something
type Announcer struct {
	actor.Actor
}

// ConsumeEvent for notifiers
func (a *Announcer) ConsumeEvent(event *events.Event) {
	a.Stream <- event
}

func NewAnnouncer(gs actor.StreamInterface) *Announcer {
	a := actor.NewActor("Announcer", gs)
	announcer := new(Announcer)
	announcer.Actor = a
	announcer.Actor.ProcessEvent = announcer.ProcessEvent
	return announcer
}

// ProcessEvent - i need print something
func (a *Announcer) ProcessEvent(event *events.Event) {
	switch event.Type {
	case events.SECOND:
		a.SendEvent("global", events.MESSAGE, "Every second, mister")
	case events.MINUTE:
		world := a.World
		a.SendEvent("global", events.MESSAGE, fmt.Sprintf("Игровое время: %v", world.GetDate().Format("Mon Jan _2 15:04:05 2006")))
	}
}
