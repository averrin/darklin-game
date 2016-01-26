package main

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

//Character - room-based actor
type Character struct {
	Actor
	Room *Area
}

func (a *Character) String() string {
	return fmt.Sprintf("{Name: %s, Room: %s}", a.Name, a.Room.Name)
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
	db.C("npc").Upsert(bson.M{"name": a.Name}, a.State)
}

//ChangeRoom - enter to new room
func (a *NPC) ChangeRoom(room *Area) {
	a.BroadcastRoom(ROOMEXIT, "Покинул комнату", a.Name, a.Room)
	delete(a.Room.Streams, a.Name)
	delete(a.Room.NPCs, a.Name)
	a.Streams["room"] = &room.Stream
	a.Room = room
	a.State.Room = room.Name
	go a.UpdateState()
	room.Streams[a.Name] = &a.Stream
	room.NPCs[a.Name] = a
	a.BroadcastRoom(ROOMENTER, "Вошел в комнату", a.Name, a.Room)
	a.SendEvent("room", ROOMENTER, nil)
	a.Stream <- NewEvent(ROOMCHANGED, a.Room.Name, "global")
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
