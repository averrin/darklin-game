package rooms

import (
	"actor"
	"objects"
)

//ShopInit - room init
func (a *Room) ShopInit() {
	t := objects.NewTable()
	a.AddObject("Table", &t)
	// if a.State.New {
	item, _ := a.World.GetItem("Old key")
	t.AddItem(item)
	// }
}

//NewShop - constructor
func NewShop(gs actor.StreamInterface) *Room {
	hall := NewRoom("Shop", "Магазин. В центре комнаты стоит стол [Table].", (*Room).ShopInit, []string{"Hall"}, gs)
	world := hall.World

	// hall.Handlers[events.LIGHT] = hall.HallLight

	world.AddRoom("Shop", &hall)
	// go hall.Live()

	return &hall
}
