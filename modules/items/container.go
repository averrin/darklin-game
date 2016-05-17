package items

import (
	actor "../actor"
	"bytes"
	"html/template"
)

//Container - items holder
type Container struct {
	Items map[string]actor.ItemInterface
}

//NewContainer - constructor
func NewContainer() *Container {
	container := new(Container)
	container.Items = make(map[string]actor.ItemInterface)
	return container
}

//Count -
func (a *Container) Count() int {
	return len(a.Items)
}

//AddItem -
func (a *Container) AddItem(name string, item actor.ItemInterface) {
	a.Items[name] = item
}

//GetItem -
func (a *Container) GetItem(name string) (actor.ItemInterface, bool) {
	item, ok := a.Items[name]
	return item, ok
}

//GetItems -
func (a *Container) GetItems() map[string]actor.ItemInterface {
	return a.Items
}

//RemoveItem -
func (a *Container) RemoveItem(name string) {
	delete(a.Items, name)
}

//String -
func (a *Container) String() string {
	tplString := "{{range $key, $item := .Items}}\n  * {{$item.GetDesc}} [{{$key}}]{{end}}"
	tpl, err := template.New("container").Parse(tplString)
	if err != nil {
		panic(err)
	}

	buffer := bytes.NewBuffer([]byte{})
	err = tpl.Execute(buffer, a)
	if err != nil {
		panic(err)
	}
	return buffer.String()
}
