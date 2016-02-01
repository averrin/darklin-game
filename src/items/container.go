package items

import (
	"actor"
	"bytes"
	"text/template"
)

type Container struct {
	Items map[string]actor.ItemInterface
}

func NewContainer() *Container {
	container := new(Container)
	container.Items = make(map[string]actor.ItemInterface)
	return container
}

func (a *Container) Count() int {
	return len(a.Items)
}

func (a *Container) AddItem(name string, item actor.ItemInterface) {
	a.Items[name] = item
}

func (a *Container) GetItem(name string) (actor.ItemInterface, bool) {
	item, ok := a.Items[name]
	return item, ok
}

func (a *Container) GetItems() map[string]actor.ItemInterface {
	return a.Items
}

func (a *Container) RemoveItem(name string) {
	delete(a.Items, name)
}

func (a *Container) String() string {
	tplString := "{{range $key, $item := .Items}}\t{{$key}} -- {{$item.GetDesc}}\n{{end}}"
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
