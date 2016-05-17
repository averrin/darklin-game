package items

import (
	actor "../actor"
	"fmt"
)

//Item -
type Item struct {
	Name string
	Desc string
}

//GetName -
func (a *Item) GetName() string {
	return a.Name
}

//GetDesc -
func (a *Item) GetDesc() string {
	return a.Desc
}

//String -
func (a *Item) String() string {
	return fmt.Sprintf("{Name: %s, Desc: %s}", a.Name, a.Desc)
}

func (a *Item) Inspect() string {
	return a.Desc
}

func (a *Item) Use(item actor.ItemInterface) interface{} {
	return "И ничего не произошло."
}
