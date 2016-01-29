package rooms

import (
	"actor"
	"log"
)
import "items"

//ShopInit - room init
func (a *Room) ShopInit() {
	key := new(items.Item)
	key.Name = "Key"
	key.Desc = "Огромный старый ключ."
	a.Items["Key"] = key
	log.Println(a.Items)
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
