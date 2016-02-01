package rooms

import "actor"

//ShopInit - room init
func (a *Room) ShopInit() {
	if a.State.New {
		item, _ := a.World.GetItem("Key")
		a.AddItem(item)
	}
}

//NewShop - constructor
func NewShop(gs actor.StreamInterface) *Room {
	hall := NewRoom("Shop", "Магазин.", (*Room).ShopInit, []string{"Hall"}, gs)
	world := hall.World

	// hall.Handlers[events.LIGHT] = hall.HallLight

	world.AddRoom("Shop", &hall)
	// go hall.Live()

	return &hall
}
