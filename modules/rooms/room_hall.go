package rooms

import (
	actor "../actor"
	events "../events"
	npc "../npc"
	objects "../objects"
)

//NewHall - constructor
func NewHall(gs actor.StreamInterface) *Room {
	hall := NewRoom("Hall", "Это холл. Большая, светлая комната. У стены стоит сундук [Chest].", (*Room).HallInit, []string{"Store", "Shop"}, gs)
	world := hall.World

	hall.Handlers[events.LIGHT] = hall.HallLight

	world.AddRoom("Hall", &hall)
	// go hall.Live()

	return &hall
}

//HallInit -
func (a *Room) HallInit() {
	world := a.World
	gs := *world.GetGlobal()
	mik := npc.NewMik(gs)
	// log.Println(mik.State.New)
	t := objects.NewChest()
	a.AddObject("Chest", &t)
	if mik.State.New {
		a.AddNPC(mik)
	}
	item, _ := a.World.GetItem("Sword")
	t.AddItem(item)
	t.Lock("Old key")
	go a.UpdateState()
	go mik.Live()
}

//HallLight -
func (a *Room) HallLight(event *events.Event) bool {
	if !event.Payload.(bool) {
		a.BroadcastRoom(events.SYSTEMMESSAGE, "Стало как-то неуютно", a.Name)
	}
	return false
}
