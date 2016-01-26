package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/websocket"
	"github.com/ugorji/go/codec"
)

// Player just someone who do something
type Player struct {
	Character
	Connection *websocket.Conn
	State      PlayerState
	Loggedin   bool
}

// ConsumeEvent of cause
func (a *Character) ConsumeEvent(event *Event) {
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
	actor.Actor.ProcessEvent = actor.ProcessEvent
	actor.Actor.ProcessCommand = actor.ProcessCommand
	return actor
}

//HashPassword - hash. password.
func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr
}

//Login user
func (a *Player) Login(login string, password string) (string, bool) {
	// delete(a.Streams, p.Name)
	password = HashPassword(password)
	a.State = *new(PlayerState)
	s := a.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	n, _ := db.C("players").Find(bson.M{"name": login}).Count()
	if n != 0 {
		db.C("players").Find(bson.M{"name": login}).One(&a.State)
		if password != a.State.Password {
			return "Неверный пароль", false
		}
		a.State.New = false
	} else {
		a.State.New = true
		a.State.Name = login
		a.State.Password = password
		a.State.Room = "Hall"
		a.State.HP = 10
	}
	a.Name = login
	a.Loggedin = true
	go a.Live()
	a.ChangeRoom(WORLD.Rooms[a.State.Room])
	db.C("players").Upsert(bson.M{"name": a.Name}, a.State)
	// log.Println("success login", blue(tokens[1]))
	go a.Message(NewEvent(LOGGEDIN, "Вы вошли как: "+a.Name, "global"))
	return "", a.State.New
}

//UpdateState - save state into db
func (a *Player) UpdateState() {
	s := a.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	db.C("players").Update(bson.M{"name": a.Name}, a.State)
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
		a.BroadcastRoom(ROOMEXIT, "Покинул комнату", a.Name, a.Room)
		delete(a.Room.Streams, a.Name)
		delete(a.Room.Players, a)
	}
	a.Streams["room"] = &room.Stream
	a.Room = room
	a.State.Room = room.Name
	go a.UpdateState()
	room.Players[a] = a.Connection
	room.Streams[a.Name] = &a.Stream
	a.BroadcastRoom(ROOMENTER, "Вошел в комнату", a.Name, a.Room)
	a.SendEvent("room", ROOMENTER, nil)
	if prevRoom != nil {
		a.Stream <- NewEvent(ROOMCHANGED, fmt.Sprintf("Вы здесь: %v", a.Room.Name), "global")
	}
}

//ProcessEvent - event handler
func (a *Player) ProcessEvent(event *Event) {
	// log.Println(event)
	switch event.Type {
	case LOGIN:
		err, _ := a.Login(event.Payload.([]string)[1], event.Payload.([]string)[2])
		if err != "" {
			a.Message(NewEvent(ERROR, err, "global"))
		} else {
			a.SendEvent("global", LOGGEDIN, nil)
			// a.Streams[p.Name] = &p.Stream
		}
	default:
		if a.Connection != nil {
			a.Message(event)
		}
	}
}

//ProcessCommand - command handler
func (a *Player) ProcessCommand(event *Event) {}

//PlayerState - db saved state
type PlayerState struct {
	CharState
	Password string
}
