package world

import (
	"actor"
	"events"
	"npc"
	"rooms"
	"time"
	"timeStream"
)

//World - global container
type World struct {
	Rooms  map[string]actor.RoomInterface
	Global *actor.StreamInterface
	Time   actor.TimeInterface
}

//NewWorld - constructor
func NewWorld(gs actor.StreamInterface) *World {
	world := new(World)
	// gs := *gsl
	gs.SetWorld(world)
	world.Global = &gs
	// log.Println((*world.Global).GetWorld())
	world.Rooms = make(map[string]actor.RoomInterface)
	world.Time = timeStream.NewTimeStream(gs, gs.GetDate())
	go world.Time.Live()

	return world
}

//Init - create rooms
func (w *World) Init() {
	gs := *w.Global
	hall := rooms.NewHall(gs)
	go hall.Live()
	store := rooms.NewStore(gs)
	go store.Live()
	shop := rooms.NewShop(gs)
	go shop.Live()
	announcer := npc.NewAnnouncer(gs)
	go announcer.Live()
	store.Init(store)
	hall.Init(hall)
	shop.Init(hall)

	// gs.Subscribe(SECOND, announcer)
	gs.Subscribe(events.MINUTE, &announcer.Actor)
}

//AddRoom -
func (w *World) AddRoom(name string, room actor.RoomInterface) {
	w.Rooms[name] = room
}

//GetDate -
func (w *World) GetDate() time.Time {
	return w.Time.GetDate()
}

//GetGlobal -
func (w *World) GetGlobal() *actor.StreamInterface {
	return w.Global
}

//GetTime -
func (w *World) GetTime() *actor.TimeInterface {
	return &w.Time
}

//GetRoom -
func (w *World) GetRoom(name string) (*actor.RoomInterface, bool) {
	// log.Println(w.Rooms)
	room, ok := w.Rooms[name]
	return &room, ok
}
