package area

import (
	"actor"
	"core"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/websocket"
)

//Area - room for players
type Area struct {
	actor.Actor
	Players   map[actor.PlayerInterface]*websocket.Conn
	Formatter core.Formatter
	State     actor.AreaState
}

//area.NewArea constructor
func NewArea(name string, gs actor.StreamInterface) Area {
	a := actor.NewActor(name, gs)
	area := new(Area)
	area.Actor = a
	area.Players = make(map[actor.PlayerInterface]*websocket.Conn)
	formatter := core.NewFormatter()
	area.Formatter = formatter
	// area.Actor.ProcessEvent = area.ProcessEvent
	s := area.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	n, _ := db.C("rooms").Find(bson.M{"name": area.Name}).Count()
	area.State = *new(actor.AreaState)
	area.State.New = true
	area.State.Light = true
	area.State.Name = area.Name
	if n != 0 {
		db.C("rooms").Find(bson.M{"name": area.Name}).One(&area.State)
		area.State.New = false
	}
	return *area
}

//UpdateState - save state into db
func (a *Area) UpdateState() {
	s := a.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	a.State.New = false
	db.C("rooms").Upsert(bson.M{"name": a.Name}, a.State)
}

//GetPlayer by name
func (a *Area) GetPlayer(name string) *actor.PlayerInterface {
	for v := range a.Players {
		if v.GetName() == name {
			return &v
		}
	}
	return new(actor.PlayerInterface)
}
