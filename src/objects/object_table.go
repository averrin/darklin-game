package objects

import "items"

//NewTable - constructor
func NewTable() Object {
	table := new(Object)
	table.Name = "Table"
	table.Desc = "Стол. Не на что тут смотреть."
	table.State = new(ObjectState)
	container := items.NewContainer()
	table.Items = container
	return *table
}
