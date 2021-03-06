package npc

import (
	"fmt"

	actor "../actor"
	events "../events"

	"gopkg.in/mgo.v2/bson"
)

//Character - room-based actor
type Character struct {
	actor.Actor
	Desc  string
	Room  *actor.RoomInterface
	World actor.WorldInterface
}

func (a *Character) String() string {
	room := *a.Room
	return fmt.Sprintf("{Name: %s, Room: %s}", a.Name, room.GetName())
}

//NPC - just NPC
type NPC struct {
	Character
	State actor.CharState
}

//UpdateState - save state into db
func (a *NPC) UpdateState() {
	s := a.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	db.C("npc").Upsert(bson.M{"name": a.Name}, a.State)
}

//ChangeRoom - enter to new room
func (a *NPC) ChangeRoom(room *actor.RoomInterface) {
	prevRoom := *a.Room
	prevRoom.BroadcastRoom(events.ROOMEXIT, "Покинул комнату", a.Name)
	prevRoom.RemoveNPC(a.Name)
	// delete(a.Room.NPCs, a.Name)
	// room.AddNPC(a.(*actor.NPCInterface))
	// n := actor.NPCInterface(a)
	(*room).AddNPC(a)
	// room.Streams[a.Name] = &a.Stream
	// room.NPCs[a.Name] = a
	(*room).BroadcastRoom(events.ROOMENTER, "Вошел в комнату", a.Name)
	a.SendEvent("room", events.ROOMENTER, nil)
	a.Stream <- events.NewEvent(events.ROOMCHANGED, (*room).GetName(), "global")
}

// NewNPC constructor
func NewNPC(name string, desc string, gs actor.StreamInterface) NPC {
	a := actor.NewActor(name, gs)
	char := new(NPC)
	char.Desc = desc
	char.Actor = a
	char.World = gs.GetWorld()
	// formatter := NewFormatter()
	// actor.Formatter = formatter
	char.Actor.ProcessEvent = char.ProcessEvent
	s := char.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	n, _ := db.C("npc").Find(bson.M{"name": char.Name}).Count()
	char.State = *new(actor.CharState)
	char.State.New = true
	char.State.Name = char.Name
	if n != 0 {
		db.C("npc").Find(bson.M{"name": char.Name}).One(&char.State)
		room, _ := a.World.GetRoom(char.State.Room)
		(*room).AddNPC(char)
		char.State.New = false
	}
	return *char
}

//ProcessEvent from user or cmd
func (a *NPC) ProcessEvent(event *events.Event) {
	// formatter := a.Formatter
	// blue := formatter.Blue
	// yellow := formatter.Yellow
	// log.Println(a.Name, event)
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

//SetRoom -
func (a *NPC) SetRoom(room actor.RoomInterface) {
	a.Room = &room
	a.State.Room = room.GetName()
	go a.UpdateState()
}

//GetDesc -
func (a *NPC) GetDesc() string {
	return a.Desc
}

func (a *NPC) Inspect() string {
	return a.Desc
}
