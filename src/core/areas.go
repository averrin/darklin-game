package main

import (
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

//Area - room for players
type Area struct {
	Actor
	Players   map[*Player]*websocket.Conn
	Formatter Formatter
}

// NewArea constructor
func NewArea(name string, gs *chan *Event) *Area {
	a := NewActor(name, gs)
	actor := new(Area)
	actor.Actor = *a
	actor.Players = make(map[*Player]*websocket.Conn)
	formatter := NewFormatter()
	actor.Formatter = formatter
	actor.Actor.ProcessEvent = actor.ProcessEvent
	// s := actor.Storage.Session.Copy()
	// defer s.Close()
	// db := s.DB("darklin")
	// n, _ := db.C("state").Count()
	// actor.State = *new(GlobalState)
	// actor.State.Date = time.Date(774, 1, 1, 12, 0, 0, 0, time.UTC)
	// actor.State.New = true
	// if n != 0 {
	// 	db.C("state").Find(bson.M{}).One(&actor.State)
	// 	actor.State.New = false
	// }
	return actor
}

//ProcessEvent from user or cmd
func (a *Area) ProcessEvent(event *Event) {
	// formatter := a.Formatter
	// blue := formatter.Blue
	tokens := strings.Split(event.Payload.(string), " ")
	// log.Println(tokens, len(tokens))
	command := strings.ToLower(tokens[0])
	// log.Println(command)
	_, ok := a.Streams[event.Sender]
	log.Println("Recv command " + command + " from " + event.Sender)
	if ok == false && command != "login" && event.Sender != "cmd" {
		log.Println("Discard command " + command + " from " + event.Sender)
		return
	}
	switch command {
	default:
		log.Println(a.Name, "forward", event)
		a.ForwardEvent("global", event)
	}
}
