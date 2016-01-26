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
