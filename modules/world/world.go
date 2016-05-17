package world

import (
	"io/ioutil"
	"path/filepath"
	"time"

	actor "../actor"
	events "../events"
	items "../items"
	npc "../npc"
	rooms "../rooms"
	timeStream "../timeStream"

	yaml "gopkg.in/yaml.v2"
)

//World - global container
type World struct {
	Rooms  map[string]actor.RoomInterface
	Global *actor.StreamInterface
	Time   actor.TimeInterface
	Items  actor.ItemContainerInterface
}

//NewWorld - constructor
func NewWorld(gs actor.StreamInterface) *World {
	world := new(World)
	// gs := *gsl
	gs.SetWorld(world)
	world.Global = &gs
	// log.Println((*world.Global).GetWorld())
	world.Rooms = make(map[string]actor.RoomInterface)
	// world.Items = make(map[string]actor.ItemInterface)
	container := items.NewContainer()
	world.Items = container
	world.Time = timeStream.NewTimeStream(gs, gs.GetDate())
	go world.Time.Live()

	return world
}

//Init - create rooms
func (w *World) Init() {

	rooms.InitHandlers()
	allItems := ReadItems()
	for _, itemDesc := range allItems {
		item := new(items.Item)
		item.Name = itemDesc.Name
		item.Desc = itemDesc.Desc
		w.Items.AddItem(item.Name, item)
	}

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
	shop.Init(shop)

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

//AddItem -
func (w *World) AddItem(item actor.ItemInterface) {
	w.Items.AddItem(item.GetName(), item)
}

//RemoveItem -
func (w *World) RemoveItem(name string) {
	w.Items.RemoveItem(name)
}

//GetItem -
func (w *World) GetItem(name string) (actor.ItemInterface, bool) {
	return w.Items.GetItem(name)
}

type ItemDesc struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
	Desc string `yaml:"desc"`
}

func ReadItems() []ItemDesc {
	filename, _ := filepath.Abs("modules/items/items_db.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var allItems []ItemDesc

	err = yaml.Unmarshal(yamlFile, &allItems)
	if err != nil {
		panic(err)
	}

	return allItems
}
