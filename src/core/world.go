package main

import "log"

type World struct {
	Rooms  map[string]*Area
	Global *GlobalStream
	Time   *TimeStream
}

func NewWorld(gs *GlobalStream) *World {
	world := new(World)
	world.Global = gs
	world.Rooms = make(map[string]*Area)
	world.Time = NewTimeStream(&gs.Stream, gs.State.Date)
	go world.Time.Live()

	return world
}

func (w *World) Init() {
	gs := w.Global
	room2 := NewArea("second", &gs.Stream)
	go room2.Live()
	w.Rooms["second"] = &room2
	hall := NewHall(&gs.Stream)
	log.Println(hall)
	hall.Init()
	// room := NewArea("Hall", &gs.Stream)
	// go room.Live()
	// world.Rooms["Hall"] = &hall
}

var WORLD *World
