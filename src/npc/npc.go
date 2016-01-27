package npc

import (
	"actor"
	"events"
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

//Character - room-based actor
type Character struct {
	actor.Actor
	Room  *actor.RoomInterface
	World *actor.WorldInterface
}

func (a *Character) String() string {
	room := *a.Room
	return fmt.Sprintf("{Name: %s, Room: %s}", a.Name, room.GetName())
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
func (a *NPC) ChangeRoom(rooml *actor.RoomInterface) {
	room := *rooml
	prevRoom := *a.Room
	prevRoom.BroadcastRoom(events.ROOMEXIT, "Покинул комнату", a.Name)
	prevRoom.RemoveNPC(a.Name)
	// delete(a.Room.NPCs, a.Name)
	a.Streams["room"] = &room.Stream
	a.Room = room
	a.State.Room = room.Name
	go a.UpdateState()
	room.Streams[a.Name] = &a.Stream
	room.NPCs[a.Name] = a
	room.BroadcastRoom(events.ROOMENTER, "Вошел в комнату", a.Name)
	a.SendEvent("room", events.ROOMENTER, nil)
	a.Stream <- NewEvent(events.ROOMCHANGED, a.Room.Name, "global")
}

// NewNPC constructor
func NewNPC(name string, gs *chan *events.Event, roomName string) NPC {
	a := NewActor(name, gs)
	char := new(NPC)
	char.Actor = *a
	// formatter := NewFormatter()
	// actor.Formatter = formatter
	char.Actor.ProcessEvent = char.ProcessEvent
	s := char.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	n, _ := db.C("npc").Find(bson.M{"name": char.Name}).Count()
	char.State = *new(CharState)
	char.State.New = true
	char.State.Name = char.Name
	room := a.World.GetRoom(roomName)
	char.State.Room = room.Name
	char.Room = room
	if n != 0 {
		db.C("npc").Find(bson.M{"name": char.Name}).One(&char.State)
		char.State.New = false
		char.Room = char.World.GetRoom(char.State.Room)
	}
	char.Streams["room"] = &char.Room.Stream
	char.Room.Streams[char.Name] = &char.Stream
	return *char
}

//ProcessEvent from user or cmd
func (a *NPC) ProcessEvent(event *Event) {
	// formatter := a.Formatter
	// blue := formatter.Blue
	// yellow := formatter.Yellow
	handler, ok := a.Handlers[event.Type]
	switch event.Type {
	case events.COMMAND:
		if ok {
			handled := handler(event)
			if handled {
				return
			}
		}
		a.ProcessCommand(event)
	default:
		if ok {
			_ = handler(event)
		}
	}
}
