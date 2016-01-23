package main

import (
	"log"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/websocket"
)

//Area - room for players
type Area struct {
	Actor
	Players   map[*Player]*websocket.Conn
	Formatter Formatter
	State     AreaState
}

// NewArea constructor
func NewArea(name string, gs *chan *Event) Area {
	a := NewActor(name, gs)
	actor := new(Area)
	actor.Actor = *a
	actor.Players = make(map[*Player]*websocket.Conn)
	formatter := NewFormatter()
	actor.Formatter = formatter
	actor.Actor.ProcessEvent = actor.ProcessEvent
	s := actor.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	n, _ := db.C("rooms").Find(bson.M{"name": actor.Name}).Count()
	actor.State = *new(AreaState)
	actor.State.New = true
	actor.State.Light = true
	actor.State.Name = actor.Name
	if n != 0 {
		db.C("rooms").Find(bson.M{"name": actor.Name}).One(&actor.State)
		actor.State.New = false
	}
	return *actor
}

//ProcessEvent from user or cmd
func (a *Area) ProcessEvent(event *Event) {
	// formatter := a.Formatter
	// blue := formatter.Blue
	// yellow := formatter.Yellow
	switch event.Type {
	case ROOMENTER:
		if !a.State.Light {
			a.SendEvent(event.Sender, SYSTEMMESSAGE, "В комнате темно")
		}
		log.Println(a.Name, event)
	case COMMAND:
		a.ProcessCommand(event)
	}
}

//ProcessCommand from user or cmd
func (a *Area) ProcessCommand(event *Event) {
	// formatter := a.Formatter
	// blue := formatter.Blue
	tokens := strings.Split(event.Payload.(string), " ")
	// log.Println(tokens, len(tokens))
	command := strings.ToLower(tokens[0])
	// log.Println(command)
	_, ok := a.Streams[event.Sender]
	log.Println(a.Name + ": Recv command '" + event.Payload.(string) + "' from " + event.Sender)
	if ok == false && command != "login" && event.Sender != "cmd" {
		log.Println("Discard command " + command + " from " + event.Sender)
		return
	}
	switch command {
	case "light":
		if len(tokens) == 2 && (tokens[1] == "on" || tokens[1] == "off") {
			if tokens[1] == "on" {
				if a.State.Light {
					go a.Broadcast(SYSTEMMESSAGE, "В комнате уже светло", a.Name)
					return
				}
				a.State.Light = true
				go a.Broadcast(SYSTEMMESSAGE, "В комнате зажегся свет", a.Name)
			} else {
				if !a.State.Light {
					go a.Broadcast(SYSTEMMESSAGE, "В комнате уже темно", a.Name)
					return
				}
				a.State.Light = false
				go a.Broadcast(SYSTEMMESSAGE, "В комнате погас свет", a.Name)
			}
			go a.UpdateState()
		}
	default:
		if strings.HasPrefix(command, "/") {
			a.Broadcast(MESSAGE, event.Payload.(string)[1:len(event.Payload.(string))], event.Sender)
		} else {
			log.Println(a.Name, "forward", event)
			a.ForwardEvent("global", event)
		}
	}
}

//UpdateState - save state into db
func (a *Area) UpdateState() {
	s := a.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	a.State.New = false
	db.C("rooms").Upsert(bson.M{"name": a.Name}, a.State)
}

//AreaState - db saved state
type AreaState struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string

	Light bool

	New bool
}
