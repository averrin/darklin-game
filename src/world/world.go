package world

import (
	"actor"
	"area"
	"core"
	"events"
	"log"
	"rooms"
	"time"
	"timeStream"
)

type World struct {
	Rooms  map[string]*area.Area
	Global *actor.StreamInterface
	Time   *timeStream.TimeStream
}

func NewWorld(gsl *actor.StreamInterface) *World {
	world := new(World)
	gs := *gsl
	world.Global = gsl
	wi := actor.WorldInterface(world)
	gs.SetWorld(&wi)
	world.Rooms = make(map[string]*area.Area)
	world.Time = timeStream.NewTimeStream(gs.GetStream(), gs.GetDate())
	go world.Time.Live()

	return world
}

func (w *World) Init() {
	gs := *w.Global
	room2 := area.NewArea("second", gs.GetStream())
	room2.Desc = "Абстрактная комната, не имеющая индивидуальности."
	go room2.Live()
	w.Rooms["second"] = room2
	hall := rooms.NewHall(&gs.Stream)
	hall.Init()
	announcer := core.NewAnnouncer(&gs.Stream)
	go announcer.Live()

	// gs.Subscribe(SECOND, announcer)
	gs.Subscribe(events.MINUTE, &announcer.Actor)
}

func (w *World) AddRoom(name string, room *actor.RoomInterface) {
	log.Fatal("not implemented")
}

func (w *World) GetDate() time.Time {
	return w.Time.Date
}

func (w *World) GetGlobal() *actor.StreamInterface {
	return w.Global
}

func (w *World) GetRoom(name string) (*actor.RoomInterface, bool) {
	room, ok := w.Rooms[name]
	return room, ok
}
