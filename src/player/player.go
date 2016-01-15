package player

import (
	"actor"
	"events"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Player just someone who do something
type Player struct {
	actor.Actor
	Connection *websocket.Conn
}

// ConsumeEvent of couse
func (a Player) ConsumeEvent(event events.Event) {
	a.Stream <- event
}

// NewPlayer because i, sucj in golang yet
func NewPlayer(name string, gs chan events.Event) *Player {
	log.Println("New player: ", name)
	a := actor.NewActor(name, gs)
	actor := new(Player)
	actor.Actor = *a
	return actor
}

// Live - i need print something
func (a Player) Live() {
	log.Println("Player", a.Name, "Live")
	for {
		event := <-a.Stream
		log.Println(event)
		a.NotifySubscribers(event)
		msg := fmt.Sprintf("%v: %v", event.Sender, event.Payload)
		log.Println(msg)
		err := a.Connection.WriteMessage(websocket.TextMessage, []byte(msg))
		log.Println(err)
	}
}
