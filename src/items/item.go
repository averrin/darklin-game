package items

import "fmt"

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

//GetName -
func (a *Item) String() string {
	return fmt.Sprintf("{Name: %s, Desc: %s}", a.Name, a.Desc)
}