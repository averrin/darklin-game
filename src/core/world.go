package main

type World struct {
	Rooms map[string]*Area
}

func NewWorld(gs *GlobalStream) *World {
	world := new(World)
	world.Rooms = make(map[string]*Area)
	room := NewArea("first", &gs.Stream)
	go room.Live()
	world.Rooms["first"] = &room
	room2 := NewArea("second", &gs.Stream)
	go room2.Live()
	world.Rooms["second"] = &room2
	return world
}

var WORLD *World
