package main

import "github.com/gorilla/websocket"

//Area - room for players
type Area struct {
	Actor
	Players   map[*Player]*websocket.Conn
	Storage   *Storage
	Formatter Formatter
}

// NewArea constructor
func NewArea(name string, gs chan *Event) *Area {
	a := NewActor(name, gs)
	actor := new(Area)
	actor.Actor = *a
	actor.Players = make(map[*Player]*websocket.Conn)
	actor.Storage = NewStorage()
	formatter := NewFormatter()
	actor.Formatter = formatter
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
