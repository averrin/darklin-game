package rooms

import (
	"actor"
	"events"
	"npc"
)

func NewHall(gs actor.StreamInterface) *Room {
	hall := NewRoom("Hall", "Это холл. Большая, светлая комната.", (*Room).HallInit, gs)
	world := hall.World

	hall.Handlers[events.LIGHT] = hall.HallLight

	world.AddRoom("Hall", &hall)
	// go hall.Live()

	return &hall
}

func (a *Room) HallInit() {
	world := a.World
	gs := *world.GetGlobal()
	mik := npc.NewMik(gs)
	// log.Println(mik.State.New)
	if mik.State.New {
		a.AddNPC(mik)
	}
	go mik.Live()
}

func (a *Room) HallLight(event *events.Event) bool {
	if !event.Payload.(bool) {
		a.BroadcastRoom(events.SYSTEMMESSAGE, "Стало как-то неуютно", a.Name)
	}
	return false
}
