package rooms

import (
	"actor"
	"area"
	"events"
	"fmt"
	"log"
)

type Room struct {
	area.Area
	NPCs map[string]actor.NPCInterface
}

func NewRoom(name string, gs *chan *events.Event) *Room {
	a := area.NewArea(name, gs)
	room := new(Room)
	room.Area = *a
	room.NPCs = make(map[string]actor.NPCInterface)
	return room
}

func (a *Room) String() string {
	return fmt.Sprintf("{Name: %s, Players: %d, NPCs: %d}", a.Name, len(a.Players), len(a.NPCs))
}

func (a *Room) AddNPC(actor.NPCInterface) {
	log.Fatal("not implemented")
}

func (a *Room) RemoveNPC(name string) {
	log.Fatal("not implemented")
}

func (a *Room) AddPlayer(actor.PlayerInterface) {
	log.Fatal("not implemented")
}

// BroadcastRoom - send all
func (a *Room) BroadcastRoom(eventType events.EventType, payload interface{}, sender string) {
	event := events.NewEvent(eventType, payload, sender)
	defer func() { recover() }()
	for v := range a.Players {
		p := *v
		if p.GetName() == sender {
			continue
		}
		stream := *p.GetStream()
		stream <- event
	}
	for name, npc := range a.NPCs {
		if name == sender {
			continue
		}
		stream := *npc.GetStream()
		stream <- event
	}
}

func (a *Room) GetState() actor.AreaState {
	return a.State
}
