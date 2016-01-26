package rooms

import (
	"area"
	"events"
	"npc"
	"world"
)

func NewHall(gs *chan *events.Event) *area.Area {
	hall := area.NewArea("Hall", gs)
	world.WORLD.Rooms["Hall"] = hall
	hall.Desc = "Это холл. Большая, светлая комната."

	hall.Handlers[events.LIGHT] = hall.HallLight

	go hall.Live()

	return hall
}

func (a *area.Area) Init() {
	mik := npc.NewMik(&world.WORLD.Global.Stream)
	go mik.Live()
}

func (a *area.Area) HallLight(event *events.Event) bool {
	if !event.Payload.(bool) {
		a.BroadcastRoom(events.SYSTEMMESSAGE, "Стало как-то неуютно", a.Name, a)
	}
	return false
}
