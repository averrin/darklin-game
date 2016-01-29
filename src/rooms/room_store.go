package rooms

import "actor"

func (a *Room) StoreInit() {
}

func NewStore(gs actor.StreamInterface) *Room {
	hall := NewRoom("Store", "Еще не придумал, толи магазин, толи склад.", (*Room).StoreInit, gs)
	world := hall.World

	// hall.Handlers[events.LIGHT] = hall.HallLight

	world.AddRoom("Store", &hall)
	// go hall.Live()

	return &hall
}
