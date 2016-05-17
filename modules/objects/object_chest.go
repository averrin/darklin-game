package objects

import (
	actor "../actor"
	items "../items"
)

//Chest - lockable container
type Chest struct {
	Object
	Locked bool
	Key    string
	State  *ChestState
}

//ChestState - chest state
type ChestState struct {
	Locked bool
	Items  []string
}

//NewChest - constructor
func NewChest() Chest {
	chest := new(Chest)
	chest.Name = "Chest"
	chest.Desc = "Сундук."
	chest.State = new(ChestState)
	chest.State.Locked = false
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
	a.State.Locked = true
}

//Unlock -
func (a *Chest) Unlock() {
	a.Locked = false
	a.State.Locked = false
}

//Use - use item on object
func (a *Chest) Use(item actor.ItemInterface) interface{} {
	if item.GetName() != a.Key {
		return "И ничего не произошло."
	}
	a.Unlock()
	return "Сундук открылся."
}

//GetItem -
func (a *Chest) GetItem(name string) (actor.ItemInterface, bool) {
	if !a.Locked {
		return a.Items.GetItem(name)
	}
	return nil, false
}

//GetItems -
func (a *Chest) GetItems() map[string]actor.ItemInterface {
	if !a.Locked {
		return a.Items.GetItems()
	}
	return nil
}

//GetState - return state of chest
func (a *Chest) GetState() interface{} {
	return a.State
}

//AddItem -
func (a *Chest) AddItem(item actor.ItemInterface) {
	// pos := actor.Index(a.State.Items, item.GetName())
	a.Items.AddItem(item.GetName(), item)
	// if pos == -1 {
	a.State.Items = append(a.State.Items, item.GetName())
	// 	a.UpdateState()
	// }
}

//RemoveItem -
func (a *Chest) RemoveItem(name string) {
	a.Items.RemoveItem(name)
	pos := actor.Index(a.State.Items, name)
	a.State.Items = append(a.State.Items[:pos], a.State.Items[pos+1:]...)
	// a.UpdateState()
}
