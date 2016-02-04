package objects

import (
	"actor"
	"fmt"
)

//Object - object in room
type Object struct {
	Name  string
	Desc  string
	Items actor.ItemContainerInterface
	State *ObjectState
}

type ObjectState struct {
	Items []string
}

//GetName -
func (a *Object) GetName() string {
	return a.Name
}

//GetDesc -
func (a *Object) GetDesc() string {
	return a.Desc
}

//String -
func (a *Object) String() string {
	return fmt.Sprintf("{Name: %s, Desc: %s}", a.Name, a.Desc)
}

//AddItem -
func (a *Object) AddItem(item actor.ItemInterface) {
	// pos := actor.Index(a.State.Items, item.GetName())
	a.Items.AddItem(item.GetName(), item)
	// if pos == -1 {
	a.State.Items = append(a.State.Items, item.GetName())
	// 	a.UpdateState()
	// }
}

//RemoveItem -
func (a *Object) RemoveItem(name string) {
	a.Items.RemoveItem(name)
	pos := actor.Index(a.State.Items, name)
	a.State.Items = append(a.State.Items[:pos], a.State.Items[pos+1:]...)
	// a.UpdateState()
}

//GetItem -
func (a *Object) GetItem(name string) (actor.ItemInterface, bool) {
	return a.Items.GetItem(name)
}

//GetItems -
func (a *Object) GetItems() map[string]actor.ItemInterface {
	return a.Items.GetItems()
}

func (a *Object) Inspect() string {
	r := a.Desc
	if a.Items.Count() > 0 {
		r += "\nПредметы:"
		r += a.Items.String()
	}
	return r
}

func (a *Object) Use(item actor.ItemInterface) interface{} {
	return "И ничего не произошло."
}

func (a *Object) GetState() interface{} {
	return a.State
}
