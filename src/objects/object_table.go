package objects

//NewTable - constructor
func NewTable() Object {
	table := new(Object)
	table.Name = "Table"
	table.Desc = "Стол. Не на что тут смотреть"
	return *table
}
