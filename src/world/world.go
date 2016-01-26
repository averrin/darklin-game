package world

import (
	"area"
	"core"
	"events"
	"rooms"
	"time"
	"timeStream"
)

type StreamInterface interface {
	Live()
	SetWorld(*World)
	GetStream() *chan *events.Event
	GetDate() time.Time
}

type World struct {
	Rooms  map[string]*area.Area
	Global *StreamInterface
	Time   *timeStream.TimeStream
}

func NewWorld(gsl *StreamInterface) *World {
	world := new(World)
	gs := *gsl
	world.Global = gsl
	gs.SetWorld(world)
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

var WORLD *World
