package main

import "log"

func NewHall(gs *chan *Event) Area {
	hall := NewArea("Hall", gs)
	WORLD.Rooms["Hall"] = &hall

	hall.Handlers[LIGHT] = hall.HallLight

	go hall.Live()

	return hall
}

func (a *Area) Init() {
	mik := NewMik(&WORLD.Global.Stream)
	go mik.Live()
}

func (a *Area) HallLight(event *Event) bool {
	if !event.Payload.(bool) {
		a.BroadcastRoom(SYSTEMMESSAGE, "Стало как-то неуютно", a.Name, a)
	}
	return false
}
