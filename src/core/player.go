package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/websocket"
	"github.com/ugorji/go/codec"
)

// Player just someone who do something
type Player struct {
	Actor
	Connection *websocket.Conn
	State      PlayerState
	Loggedin   bool
	Room       *Area
}

// ConsumeEvent of couse
func (a *Player) ConsumeEvent(event *Event) {
	a.Stream <- event
}

// NewPlayer because i, sucj in golang yet
func NewPlayer(name string, gs *chan *Event) *Player {
	// green := color.New(color.FgGreen).SprintFunc()
	// log.Println("New player: ", green(name))
	a := NewActor(name, gs)
	actor := new(Player)
	actor.Actor = *a
	actor.Loggedin = false
	return actor
}

//Login user
func (a *Player) Login(login string, password string) (string, bool) {
	// delete(a.Streams, p.Name)
	a.Name = login
	s := a.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	n, _ := db.C("players").Find(bson.M{"name": a.Name}).Count()
	a.State = *new(PlayerState)
	a.State.New = true
	a.State.Name = a.Name
	a.State.Password = password
	a.State.Room = "first"
	a.State.HP = 10
	if n != 0 {
		db.C("players").Find(bson.M{"name": a.Name}).One(&a.State)
		a.State.New = false
	}
	a.Loggedin = true
	a.ChangeRoom(WORLD.Rooms[a.State.Room])
	db.C("players").Upsert(bson.M{"name": a.Name}, a.State)
	go a.Live()
	// log.Println("success login", blue(tokens[1]))
	a.Message(NewEvent(LOGGEDIN, "Вы вошли как: "+a.Name, a.Name))
	return "success", true
}

//UpdateState - save state into db
func (a *Player) UpdateState() {
	s := a.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	db.C("players").Update(bson.M{"name": a.Name}, a.State)
}

// Live - i need print something
func (a *Player) Live() {
	// log.Println("Player", a.Name, "Live")
	for a.Loggedin {
		event, ok := <-a.Stream
		if !ok {
			return
		}
		a.NotifySubscribers(event)
		switch event.Type {
		case CLOSE:
			a.Loggedin = false
			break
		default:
			a.Message(event)
		}
	}
	close(a.Stream)
	log.Println("Exit from Live of " + a.Name)
}

//Message - send event direct to ws
func (a *Player) Message(event *Event) {
	var msg []byte
	// var b []byte
	var mh codec.MsgpackHandle
	enc := codec.NewEncoderBytes(&msg, &mh)
	err := enc.Encode(event)
	if err != nil {
		log.Fatal(err)
	}
	_ = a.Connection.WriteMessage(websocket.TextMessage, msg)
}

//ChangeRoom - enter to new room
func (a *Player) ChangeRoom(room *Area) {
	prevRoom := a.Room
	if prevRoom != nil {
		a.BroadcastRoom(ROOMEXIT, "Exit from room "+a.Room.Name, a.Name, a.Room)
		delete(a.Room.Streams, a.Name)
		delete(a.Room.Players, a)
	}
	a.Streams["room"] = &room.Stream
	a.Room = room
	a.State.Room = room.Name
	go a.UpdateState()
	room.Players[a] = a.Connection
	room.Streams[a.Name] = &a.Stream
	a.BroadcastRoom(ROOMENTER, "Enter into room "+a.Room.Name, a.Name, a.Room)
	if prevRoom != nil {
		a.Stream <- NewEvent(ROOMCHANGED, fmt.Sprintf("You are here: %v", a.Room.Name), "global")
	}
}

//PlayerState - db saved state
type PlayerState struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string
	Password string

	Room string
	HP   int

	New bool
}
