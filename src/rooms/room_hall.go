package rooms

import (
	"events"
	"npc"
	"world"
)

func NewHall(gs *chan *events.Event) *Area {
	hall := area.NewArea("Hall", gs)
	world.WORLD.Rooms["Hall"] = hall
	hall.Desc = "Это холл. Большая, светлая комната."

	hall.Handlers[events.LIGHT] = hall.HallLight

	go hall.Live()

	return hall
}

func (a *Area) Init() {
	mik := npc.NewMik(&world.WORLD.Global.Stream)
	go mik.Live()
}

func (a *Area) HallLight(event *events.Event) bool {
	if !event.Payload.(bool) {
		a.BroadcastRoom(events.SYSTEMMESSAGE, "Стало как-то неуютно", a.Name, a)
	}
	return false
}
