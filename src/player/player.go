package player

import (
	"actor"
	"crypto/sha256"
	"encoding/hex"
	"events"
	"fmt"
	"items"
	"log"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/websocket"
	"github.com/ugorji/go/codec"
)

// Player just someone who do something
type Player struct {
	actor.Actor
	Room       *actor.RoomInterface
	Connection *websocket.Conn
	State      State
	Loggedin   bool
	Items      actor.ItemContainerInterface
}

// ConsumeEvent of cause
func (a *Player) ConsumeEvent(event *events.Event) {
	a.Stream <- event
}

// NewPlayer because i, sucj in golang yet
func NewPlayer(name string, gs actor.StreamInterface) actor.PlayerInterface {
	// green := color.New(color.FgGreen).SprintFunc()
	// log.Println("New player: ", green(name))
	a := actor.NewActor(name, gs)
	p := new(Player)
	p.Actor = a
	p.Loggedin = false
	p.Actor.ProcessEvent = p.ProcessEvent
	p.Actor.ProcessCommand = p.ProcessCommand
	container := items.NewContainer()
	p.Items = container
	return p
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
	a.State = *new(State)
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
		for _, name := range a.State.Items {
			item, _ := a.World.GetItem(name)
			a.AddItem(item)
		}
	} else {
		a.State.New = true
		a.State.Name = login
		a.State.Password = password
		a.State.Room = "Hall"
		a.State.HP = 10
	}
	a.Name = login
	a.Loggedin = true
	go a.Message(events.NewEvent(events.LOGGEDIN, "Вы вошли как: "+a.Name, "global"))
	go a.Live()
	room, _ := a.World.GetRoom(a.State.Room)
	a.ChangeRoom(room)
	db.C("players").Upsert(bson.M{"name": a.Name}, a.State)
	// log.Println("success login", blue(tokens[1]))
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
func (a *Player) Message(event *events.Event) {
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
func (a *Player) ChangeRoom(room *actor.RoomInterface) {
	prevRoom, _ := a.GetRoom()
	if prevRoom != nil {
		(*prevRoom).BroadcastRoom(events.ROOMEXIT, "Покинул комнату", a.Name)
		(*prevRoom).RemovePlayer(a)
	}
	a.Streams["room"] = (*room).GetStream()
	// log.Println(room, &room)
	a.Room = room
	a.State.Room = (*room).GetName()
	go a.UpdateState()
	(*room).AddPlayer(a)
	(*room).BroadcastRoom(events.ROOMENTER, "Вошел в комнату", a.Name)
	a.SendEvent("room", events.ROOMENTER, nil)
	if prevRoom != nil {
		a.Stream <- events.NewEvent(events.ROOMCHANGED, fmt.Sprintf("Вы здесь: %v", (*room).GetName()), "global")
	}
}

//ProcessEvent - event handler
func (a *Player) ProcessEvent(event *events.Event) {
	// log.Println(event)
	switch event.Type {
	case events.LOGIN:
		err, _ := a.Login(event.Payload.([]string)[1], event.Payload.([]string)[2])
		if err != "" {
			a.Message(events.NewEvent(events.ERROR, err, "global"))
		} else {
			a.SendEvent("global", events.LOGGEDIN, nil)
			// a.Streams[p.Name] = &p.Stream
		}
	case events.STATUS:
		a.Message(events.NewEvent(events.STATUS, fmt.Sprintf("Предметы:\n%v", a.Items), a.Name))
	default:
		if a.Connection != nil {
			a.Message(event)
		}
	}
}

//ProcessCommand - command handler
func (a *Player) ProcessCommand(event *events.Event) {}

//GetConnection -
func (a *Player) GetConnection() *websocket.Conn {
	return a.Connection
}

//SetConnection -
func (a *Player) SetConnection(c *websocket.Conn) {
	a.Connection = c
}

//GetRoom -
func (a *Player) GetRoom() (*actor.RoomInterface, bool) {
	if a.Room != nil {
		room := a.Room
		return room, true
	}
	return nil, false
}

//AddItem -
func (a *Player) AddItem(item actor.ItemInterface) {
	pos := actor.Index(a.State.Items, item.GetName())
	a.Items.AddItem(item.GetName(), item)
	if pos == -1 {
		a.State.Items = append(a.State.Items, item.GetName())
		a.UpdateState()
	}
}

//RemoveItem -
func (a *Player) RemoveItem(name string) {
	a.Items.RemoveItem(name)
	pos := actor.Index(a.State.Items, name)
	a.State.Items = append(a.State.Items[:pos], a.State.Items[pos+1:]...)
	a.UpdateState()
}

//GetItem -
func (a *Player) GetItem(name string) (actor.ItemInterface, bool) {
	return a.Items.GetItem(name)
}

//GetItems -
func (a *Player) GetItems() map[string]actor.ItemInterface {
	return a.Items.GetItems()
}

//State - db saved state
type State struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string

	Room  string
	Items []string
	HP    int

	New      bool
	Password string
}
