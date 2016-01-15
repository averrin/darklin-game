package player

import (
	"actor"
	"events"
	"fmt"
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
	a := actor.NewActor("player", gs)
	actor := new(Player)
	actor.Actor = *a
	return actor
}

// Live - i need print something
func (a Player) Live() {
	for {
		event := <-a.Stream
		a.NotifySubscribers(event)
		log.Println(fmt.Sprintf("%v: %v", event.Sender, event.Payload))
	}
}
