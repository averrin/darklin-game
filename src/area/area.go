package area

import (
	"actor"
	"events"
	"fmt"
	"log"
	"npc"
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
	NPCs      map[string]*npc.NPC
}

func (a *Area) String() string {
	return fmt.Sprintf("{Name: %s, Players: %d, NPCs: %d}", a.Name, len(a.Players), len(a.NPCs))
}

// NewArea constructor
func NewArea(name string, gs *chan *Event) *Area {
	a := actor.NewActor(name, gs)
	area := new(Area)
	area.Actor = *a
	area.Players = make(map[*Player]*websocket.Conn)
	area.NPCs = make(map[string]*npc.NPC)
	formatter := NewFormatter()
	area.Formatter = formatter
	area.Actor.ProcessEvent = area.ProcessEvent
	s := area.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	n, _ := db.C("rooms").Find(bson.M{"name": area.Name}).Count()
	area.State = *new(AreaState)
	area.State.New = true
	area.State.Light = true
	area.State.Name = area.Name
	if n != 0 {
		db.C("rooms").Find(bson.M{"name": area.Name}).One(&area.State)
		area.State.New = false
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
