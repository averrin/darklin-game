package rooms

import "actor"

//StoreInit - room init
func (a *Room) StoreInit() {
}

//NewStore - constructor
func NewStore(gs actor.StreamInterface) *Room {
	hall := NewRoom("Store", "Еще не придумал, толи магазин, толи склад.", (*Room).StoreInit, gs)
	world := hall.World

	// hall.Handlers[events.LIGHT] = hall.HallLight

	world.AddRoom("Store", &hall)
	// go hall.Live()

	return &hall
}
