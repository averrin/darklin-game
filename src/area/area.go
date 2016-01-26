package area

import (
	"actor"
	"events"
	"fmt"
	"log"
	"player"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/websocket"
)

//Area - room for players
type Area struct {
	actor.Actor
	Players   map[*player.Player]*websocket.Conn
	Formatter Formatter
	State     AreaState
}

//area.NewArea constructor
funcarea.NewArea(name string, gs *chan *Event) *Area {
	a := actor.NewActor(name, gs)
	area := new(Area)
	rooms.Actor = *a
	rooms.Players = make(map[*Player]*websocket.Conn)
	formatter := NewFormatter()
	rooms.Formatter = formatter
	rooms.Actor.ProcessEvent = rooms.ProcessEvent
	s := rooms.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	n, _ := db.C("rooms").Find(bson.M{"name": rooms.Name}).Count()
	rooms.State = *new(AreaState)
	rooms.State.New = true
	rooms.State.Light = true
	rooms.State.Name = rooms.Name
	if n != 0 {
		db.C("rooms").Find(bson.M{"name": rooms.Name}).One(&rooms.State)
		rooms.State.New = false
	}
	return area
}

//ProcessEvent from user or cmd
func (a *Area) ProcessEvent(event *events.Event) {
	// formatter := a.Formatter
	// blue := formatter.Blue
	// yellow := formatter.Yellow
	handler, ok := a.Handlers[event.Type]
	switch event.Type {
	case events.DESCRIBE:
		a.SendEvent(event.Sender, events.DESCRIBE, a.Desc)
	case events.ROOMENTER:
		if ok {
			handled := handler(event)
			if handled {
				return
			}
		}
		if !a.State.Light {
			a.SendEvent(event.Sender, events.SYSTEMMESSAGE, "В комнате темно")
		}
		// log.Println(a.Name, event)
	case events.COMMAND:
		if ok {
			handled := handler(event)
			if handled {
				return
			}
		}
		a.ProcessCommand(event)
	default:
		if ok {
			_ = handler(event)
		}
	}
}

//ProcessCommand from user or cmd
func (a *Area) ProcessCommand(event *events.Event) {
	// formatter := a.Formatter
	// blue := formatter.Blue
	tokens := strings.Split(event.Payload.(string), " ")
	// log.Println(tokens, len(tokens))
	command := strings.ToLower(tokens[0])
	// log.Println(command)
	_, ok := a.Streams[event.Sender]
	log.Println(fmt.Sprintf("%v: Recv command %s", a.Name, event))
	if ok == false && command != "login" && event.Sender != "cmd" {
		log.Println("Discard command " + command + " from " + event.Sender)
		return
	}
	switch command {
	case "describe":
		if tokens[1] == "room" {
			a.SendEvent(event.Sender, events.DESCRIBE, a.Desc)
		}
	case "light":
		if len(tokens) == 2 && (tokens[1] == "on" || tokens[1] == "off") {
			if tokens[1] == "on" {
				if a.State.Light {
					go a.SendEvent(event.Sender, events.SYSTEMMESSAGE, "В комнате уже светло")
					return
				}
				a.State.Light = true
				go a.Broadcast(events.SYSTEMMESSAGE, "В комнате зажегся свет", a.Name)
			} else {
				if !a.State.Light {
					go a.SendEvent(event.Sender, events.SYSTEMMESSAGE, "В комнате уже темно")
					return
				}
				a.State.Light = false
				go a.Broadcast(events.SYSTEMMESSAGE, "В комнате погас свет", a.Name)
			}
			go func() { a.Stream <- NewEvent(events.LIGHT, a.State.Light, event.Sender) }()
			go a.Broadcast(events.LIGHT, a.State.Light, a.Name)
			go a.UpdateState()
		}
	default:
		if strings.HasPrefix(command, "/") {
			a.Broadcast(events.MESSAGE, event.Payload.(string)[1:len(event.Payload.(string))], event.Sender)
		} else {
			// log.Println(a.Name, "forward", event)
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

//GetPlayer by name
func (a *Area) GetPlayer(name string) *player.Player {
	for v := range a.Players {
		if v.Name == name {
			return v
		}
	}
	return &Player{}
}

// BroadcastRoom - send all
func (a *Area) BroadcastRoom(eventType events.EventType, payload interface{}, sender string) {
	event := events.NewEvent(eventType, payload, sender)
	defer func() { recover() }()
	for p := range a.Players {
		if p.Name == sender {
			continue
		}
		p.Stream <- event
	}
	for name, npc := range a.NPCs {
		if name == sender {
			continue
		}
		npc.Stream <- event
	}
}
