package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Player just someone who do something
type Player struct {
	Actor
	Connection *websocket.Conn
	Loggedin   bool
	Room       *Area
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
		case CLOSE:
			a.Loggedin = false
			break
		default:
			a.Message(event)
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

//ChangeRoom - enter to new room
func (a *Player) ChangeRoom(room *Area) {
	prevRoom := a.Room
	if prevRoom != nil {
		a.BroadcastRoom(ROOMEXIT, "Exit from room "+a.Room.Name, a.Name, a.Room)
		delete(a.Room.Streams, a.Name)
		delete(a.Room.Players, a)
	}
	a.Streams["room"] = &room.Stream
	a.Room = room
	room.Players[a] = a.Connection
	room.Streams[a.Name] = &a.Stream
	a.BroadcastRoom(ROOMENTER, "Enter into room "+a.Room.Name, a.Name, a.Room)
	if prevRoom != nil {
		a.Stream <- NewEvent(ROOMCHANGED, fmt.Sprintf("You are here: %v", a.Room.Name), "global")
	}
}
