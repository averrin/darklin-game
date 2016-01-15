package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Player just someone who do something
type Player struct {
	Actor
	Connection *websocket.Conn
}

// ConsumeEvent of couse
func (a Player) ConsumeEvent(event Event) {
	a.Stream <- event
}

// NewPlayer because i, sucj in golang yet
func NewPlayer(name string, gs chan Event) *Player {
	log.Println("New player: ", name)
	a := NewActor(name, gs)
	actor := new(Player)
	actor.Actor = *a
	return actor
}

// Live - i need print something
func (a Player) Live() {
	log.Println("Player", a.Name, "Live")
	for {
		event := <-a.Stream
		a.NotifySubscribers(event)
		msg := fmt.Sprintf("%v: %v", event.Sender, event.Payload)
		_ = a.Connection.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}
