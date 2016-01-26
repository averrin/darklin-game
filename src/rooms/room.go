package rooms

import (
	"area"
	"fmt"
	"npc"
)

type Room struct {
	area.Area
	NPCs map[string]*npc.NPC
}

func NewRoom() *Room {
	room := area.NewArea()
	room.NPCs = make(map[string]*npc.NPC)
}

func (a *Room) String() string {
	return fmt.Sprintf("{Name: %s, Players: %d, NPCs: %d}", a.Name, len(a.Players), len(a.NPCs))
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
		npc.Stream <- event
	}
}
