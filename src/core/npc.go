package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Character - room-based actor
type Character struct {
	Actor
	Room *Area
}

//NPC - just NPC
type NPC struct {
	Character
	State CharState
}

//CharState - Basic state
type CharState struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string

	Room string
	HP   int

	New bool
}

//UpdateState - save state into db
func (a *NPC) UpdateState() {
	s := a.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	db.C("npc").Update(bson.M{"name": a.Name}, a.State)
}

//ChangeRoom - enter to new room
func (a *NPC) ChangeRoom(room *Area) {
	a.BroadcastRoom(ROOMEXIT, "Exit from room "+a.Room.Name, a.Name, a.Room)
	delete(a.Room.Streams, a.Name)
	a.Streams["room"] = &room.Stream
	a.Room = room
	a.State.Room = room.Name
	go a.UpdateState()
	room.Streams[a.Name] = &a.Stream
	a.BroadcastRoom(ROOMENTER, "Enter into room "+a.Room.Name, a.Name, a.Room)
	a.SendEvent("room", ROOMENTER, nil)
	a.Stream <- NewEvent(ROOMCHANGED, fmt.Sprintf("You are here: %v", a.Room.Name), "global")
}

// NewNPC constructor
func NewNPC(name string, gs *chan *Event, room *Area) NPC {
	a := NewActor(name, gs)
	actor := new(NPC)
	actor.Actor = *a
	// formatter := NewFormatter()
	// actor.Formatter = formatter
	actor.Actor.ProcessEvent = actor.ProcessEvent
	s := actor.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	n, _ := db.C("npc").Find(bson.M{"name": actor.Name}).Count()
	actor.State = *new(CharState)
	actor.State.New = true
	actor.State.Name = actor.Name
	actor.State.Room = room.Name
	actor.Room = room
	if n != 0 {
		db.C("npc").Find(bson.M{"name": actor.Name}).One(&actor.State)
		actor.State.New = false
		actor.Room = WORLD.Rooms[actor.State.Room]
	}
	actor.Streams["room"] = &actor.Room.Stream
	actor.Room.Streams[actor.Name] = &actor.Stream
	return *actor
}

//NewMik - nobody likes darkness
func NewMik(gs *chan *Event) NPC {
	room := WORLD.Rooms["first"]
	mik := NewNPC("Mik Rori", gs, room)
	return mik
}

//Live - Mik event loop
func (a *NPC) Live() {
	for {
		event, ok := <-a.Stream
		if !ok {
			return
		}
		a.NotifySubscribers(event)
		switch event.Type {
		case LIGHT:
			if !event.Payload.(bool) {
				a.BroadcastRoom(MESSAGE, "Эй, кто выключил свет?", a.Name, a.Room)
				a.BroadcastRoom(SYSTEMMESSAGE, "*шорох, шаги*", a.Name, a.Room)
				ne := NewEvent(COMMAND, "light on", a.Name)
				ne.ID = "Mik_light_on"
				ne.Delay = 5 * time.Second
				a.Room.Stream <- ne
			} else {
				ev, ok := a.Room.PendingEvents["Mik_light_on"]
				if ok {
					a.BroadcastRoom(MESSAGE, "То-то же!", a.Name, a.Room)
					ev.Abort = true
				}
			}
		}
	}
}
