package objects

import (
	"actor"
	"fmt"
	"items"
)

type Chest struct {
	Object
	Locked bool
	Key    string
}

//NewTable - constructor
func NewChest() Chest {
	chest := new(Chest)
	chest.Name = "Chest"
	chest.Desc = "Сундук"
	container := items.NewContainer()
	chest.Items = container
	return *chest
}

func (a *Chest) Inspect() string {
	if a.Locked {
		return "Сундук заперт."
	} else {
		return fmt.Sprintf("%s", a.Items)
	}
}

func (a *Chest) Lock(key string) {
	a.Key = key
	a.Locked = true
}

func (a *Chest) Unlock() {
	a.Locked = false
}

func (a *Chest) Use(item actor.ItemInterface) interface{} {
	if item.GetName() != a.Key {
		return "И ничего не произошло."
	} else {
		a.Unlock()
		return "Сундук открылся"
	}
}
