package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// Player just someone who do something
type Player struct {
	Actor
	Connection *websocket.Conn
	Loggedin   bool
}

// ConsumeEvent of couse
func (a *Player) ConsumeEvent(event *Event) {
	a.Stream <- event
}

// NewPlayer because i, sucj in golang yet
func NewPlayer(name string, gs *chan *Event) *Player {
	// green := color.New(color.FgGreen).SprintFunc()
	// log.Println("New player: ", green(name))
	a := NewActor(name, gs)
	actor := new(Player)
	actor.Actor = *a
	actor.Loggedin = false
	return actor
}

// Live - i need print something
func (a *Player) Live() {
	// log.Println("Player", a.Name, "Live")
	for a.Loggedin {
		event, ok := <-a.Stream
		if !ok {
			return
		}
		a.NotifySubscribers(event)
		switch event.Type {
		case HEARTBEAT:
			a.Message(event)
		case MESSAGE:
			a.Message(event)
		case LOGGEDIN:
			a.Message(event)
		case ERROR:
			a.Message(event)
		case CLOSE:
			a.Loggedin = false
			break
		}
	}
	close(a.Stream)
	log.Println("Exit from Live of " + a.Name)
}

//Message - send event direct to ws
func (a *Player) Message(event *Event) {
	msg, _ := json.Marshal(event)
	_ = a.Connection.WriteMessage(websocket.TextMessage, []byte(msg))
}
