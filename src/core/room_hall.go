package main

func NewHall(gs *chan *Event) Area {
	hall := NewArea("Hall", gs)
	WORLD.Rooms["Hall"] = &hall
	go hall.Live()

	return hall
}

func (a *Area) Init() {
	mik := NewMik(&WORLD.Global.Stream)
	go mik.Live()
}

func (a *Area) Hall(event *Event) {

}
