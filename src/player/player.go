package player

import (
	"actor"
	"events"
	"log"
)

// TestActor just someone who do something
type Player struct {
	actor.Actor
}

// ConsumeEvent of couse
func (a Player) ConsumeEvent(event events.Event) {
	a.Stream <- event
}

// NewTestActor because i, sucj in golang yet
func NewPlayer(gs chan events.Event) *Player {
	a := actor.NewActor(gs)
	actor := new(Player)
	actor.Actor = *a
	return actor
}

// Live - i need print something
func (a Player) Live() {
	for {
		event := <-a.Stream
		for _, s := range a.Subscriptions {
			if event.Type == s.Type || s.Type == events.ALL {
				go s.Subscriber.ConsumeEvent(event)
			}
		}
		log.Println(event.Payload)
	}
}
