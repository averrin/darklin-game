package world

import (
	"actor"
	"events"
	"log"
	"npc"
	"rooms"
	"time"
	"timeStream"
)

type World struct {
	Rooms  map[string]*actor.RoomInterface
	Global *actor.StreamInterface
	Time   *timeStream.TimeStream
}

func NewWorld(gs actor.StreamInterface) *World {
	world := new(World)
	// gs := *gsl
	gs.SetWorld(world)
	world.Global = &gs
	log.Println((*world.Global).GetWorld())
	world.Rooms = make(map[string]*actor.RoomInterface)
	world.Time = timeStream.NewTimeStream(gs.GetStream(), gs.GetDate())
	go world.Time.Live()

	return world
}

func (w *World) Init() {
	gs := *w.Global
	room2 := rooms.NewRoom("second", gs.GetStream())
	room2.Desc = "Абстрактная комната, не имеющая индивидуальности."
	go room2.Live()
	ri := actor.RoomInterface(room2)
	w.Rooms["second"] = &ri
	hall := rooms.NewHall(gs.GetStream())
	hall.Init()
	announcer := npc.NewAnnouncer(gs.GetStream())
	go announcer.Live()

	// gs.Subscribe(SECOND, announcer)
	gs.Subscribe(events.MINUTE, &announcer.Actor)
}

func (w *World) AddRoom(name string, room actor.RoomInterface) {
	log.Fatal("not implemented")
}

func (w *World) GetDate() time.Time {
	return w.Time.Date
}

func (w *World) GetGlobal() actor.StreamInterface {
	return *w.Global
}

func (w *World) GetRoom(name string) (actor.RoomInterface, bool) {
	room, ok := w.Rooms[name]
	return *room, ok
}
