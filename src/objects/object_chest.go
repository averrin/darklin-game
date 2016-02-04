package objects

import (
	"actor"
	"items"
)

//Chest - lockable container
type Chest struct {
	Object
	Locked bool
	Key    string
}

//NewChest - constructor
func NewChest() Chest {
	chest := new(Chest)
	chest.Name = "Chest"
	chest.Desc = "Сундук."
	container := items.NewContainer()
	chest.Items = container
	return *chest
}

//Inspect - react on lookup
func (a *Chest) Inspect() string {
	if a.Locked {
		return "Сундук заперт."
	}
	r := a.Desc
	if a.Items.Count() > 0 {
		r += "\nПредметы:"
		r += a.Items.String()
	}
	return r
}

//Lock -
func (a *Chest) Lock(key string) {
	a.Key = key
	a.Locked = true
}

//Unlock -
func (a *Chest) Unlock() {
	a.Locked = false
}

//Use - use item on object
func (a *Chest) Use(item actor.ItemInterface) interface{} {
	if item.GetName() != a.Key {
		return "И ничего не произошло."
	}
	a.Unlock()
	return "Сундук открылся."
}
