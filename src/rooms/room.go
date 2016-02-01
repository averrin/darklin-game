package rooms

import (
	"actor"
	"area"
	"events"
	"fmt"
	"items"
	"log"
	"strings"
)

//Room - room.
type Room struct {
	area.Area
	Desc    string
	NPCs    map[string]*actor.NPCInterface
	Init    func(*Room)
	ToRooms []string
	// Items   map[string]actor.ItemInterface
	Items actor.ItemContainerInterface
}

//NewRoom - constrictor
func NewRoom(name string, desc string, init func(*Room), rooms []string, gs actor.StreamInterface) Room {
	a := area.NewArea(name, gs)
	room := new(Room)
	room.Area = a
	room.Actor.ProcessEvent = room.ProcessEvent
	room.NPCs = make(map[string]*actor.NPCInterface)
	// room.Items = make(map[string]actor.ItemInterface)
	container := items.NewContainer()
	room.Items = container
	room.Desc = desc
	room.Init = init
	room.ToRooms = rooms
	for _, name := range a.State.Items {
		item, _ := room.World.GetItem(name)
		room.AddItem(item)
	}
	return *room
}

//String -
func (a *Room) String() string {
	return fmt.Sprintf("{Name: %s, Players: %d, NPCs: %d, Desc: '%s'}", a.Name, len(a.Players), len(a.NPCs), a.Desc)
}

//AddNPC -
func (a *Room) AddNPC(npc actor.NPCInterface) {
	a.Streams[npc.GetName()] = npc.GetStream()
	a.NPCs[npc.GetName()] = &npc
	npc.SetRoom(a)
	npc.SetStream("room", &a.Stream)
}

//RemoveNPC -
func (a *Room) RemoveNPC(name string) {
	delete(a.Streams, name)
	delete(a.NPCs, name)
}

//AddPlayer -
func (a *Room) AddPlayer(p actor.PlayerInterface) {
	a.Players[p] = p.GetConnection()
	a.Streams[p.GetName()] = p.GetStream()
}

//RemovePlayer -
func (a *Room) RemovePlayer(p actor.PlayerInterface) {
	delete(a.Streams, p.GetName())
	delete(a.Players, p)
}

// BroadcastRoom - send all
func (a *Room) BroadcastRoom(eventType events.EventType, payload interface{}, sender string) {
	event := events.NewEvent(eventType, payload, sender)
	defer func() { recover() }()
	for v := range a.Players {
		if v.GetName() == sender {
			continue
		}
		stream := *v.GetStream()
		stream <- event
	}
	for name, npc := range a.NPCs {
		if name == sender {
			continue
		}
		stream := *(*npc).GetStream()
		stream <- event
	}
}

//GetState -
func (a *Room) GetState() actor.AreaState {
	return *a.State
}

//ProcessEvent from user or cmd
func (a *Room) ProcessEvent(event *events.Event) {
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//ProcessCommand from user or cmd
func (a *Room) ProcessCommand(event *events.Event) {
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
	case "goto":
		go func() {
			if len(tokens) == 2 {
				p := *a.GetPlayer(event.Sender)
				w := a.World
				room, ok := w.GetRoom(tokens[1])
				if ok && stringInSlice(tokens[1], a.ToRooms) {
					p.ChangeRoom(room)
				} else {
					a.SendEvent(event.Sender, events.ERROR, fmt.Sprintf("Вы не можете перейти в эту комнату: %v", tokens[1]))
				}
			}
		}()
	case "search":
		if a.State.Light {
			a.SendEvent(event.Sender, events.DESCRIBE, fmt.Sprintf("Предметы: \n%v", a.Items))
		} else {
			go a.SendEvent(event.Sender, events.SYSTEMMESSAGE, "В комнате темно")
		}
	case "me":
		a.SendEvent(event.Sender, events.STATUS, nil)
	case "pick":
		if len(tokens) == 2 {
			p := *a.GetPlayer(event.Sender)
			item, ok := a.Items.GetItem(tokens[1])
			if ok {
				a.RemoveItem(tokens[1])
				p.AddItem(item)
				go a.SendEvent(event.Sender, events.SYSTEMMESSAGE, fmt.Sprintf("Вы подняли: %v [%v]", item.GetDesc(), item.GetName()))
			}
		}
	case "drop":
		if len(tokens) == 2 {
			p := *a.GetPlayer(event.Sender)
			item, ok := p.GetItem(tokens[1])
			if ok {
				a.AddItem(item)
				p.RemoveItem(tokens[1])
				go a.SendEvent(event.Sender, events.SYSTEMMESSAGE, fmt.Sprintf("Вы бросили: %v [%v]", item.GetDesc(), item.GetName()))
			}
		}
	case "describe":
		if len(tokens) == 2 && tokens[1] == "room" {
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
			go func() { a.Stream <- events.NewEvent(events.LIGHT, a.State.Light, event.Sender) }()
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

//GetDesc -
func (a *Room) GetDesc() string {
	return a.Desc
}

func (a *Room) AddItem(item actor.ItemInterface) {
	pos := actor.Index(a.State.Items, item.GetName())
	a.Items.AddItem(item.GetName(), item)
	if pos == -1 {
		a.State.Items = append(a.State.Items, item.GetName())
		a.UpdateState()
	}
}

func (a *Room) RemoveItem(name string) {
	a.Items.RemoveItem(name)
	pos := actor.Index(a.State.Items, name)
	a.State.Items = append(a.State.Items[:pos], a.State.Items[pos+1:]...)
	a.UpdateState()
}

func (a *Room) GetItem(name string) (actor.ItemInterface, bool) {
	return a.Items.GetItem(name)
}
