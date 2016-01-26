package world

import (
	"area"
	"core"
	"log"
	"rooms"
)

type World struct {
	Rooms  map[string]*area.Area
	Global *core.GlobalStream
	Time   *core.TimeStream
}

func NewWorld(gs *core.GlobalStream) *World {
	world := new(World)
	world.Global = gs
	world.Rooms = make(map[string]*area.Area)
	world.Time = NewTimeStream(&gs.Stream, gs.State.Date)
	go world.Time.Live()

	return world
}

func (w *World) Init() {
	gs := w.Global
	room2 := area.NewArea("second", &gs.Stream)
	room2.Desc = "Абстрактная комната, не имеющая индивидуальности."
	go room2.Live()
	w.Rooms["second"] = room2
	hall := rooms.NewHall(&gs.Stream)
	log.Println(hall)
	hall.Init()
	// room := NewArea("Hall", &gs.Stream)
	// go room.Live()
	// world.Rooms["Hall"] = &hall
}

var WORLD *World
