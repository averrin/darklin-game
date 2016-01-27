package rooms

import (
	"actor"
	"events"
	"npc"
)

func NewHall(gs actor.StreamInterface) *Room {
	hall := NewRoom("Hall", gs)
	world := hall.World
	world.AddRoom("Hall", hall)
	hall.Desc = "Это холл. Большая, светлая комната."

	hall.Handlers[events.LIGHT] = hall.HallLight

	go hall.Live()

	return hall
}

func (a *Room) Init() {
	world := a.World
	gs := *world.GetGlobal()
	mik := npc.NewMik(gs)
	go mik.Live()
}

func (a *Room) HallLight(event *events.Event) bool {
	if !event.Payload.(bool) {
		a.BroadcastRoom(events.SYSTEMMESSAGE, "Стало как-то неуютно", a.Name)
	}
	return false
}
