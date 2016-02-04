package rooms

import (
	"actor"
	"area"
	"commands"
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
	Items   actor.ItemContainerInterface
	Objects map[string]actor.ObjectInterface
}

type CommandHandler func(*Room, *events.Event, []string)

var Handlers map[commands.Command]CommandHandler

func InitHandlers() {
	Handlers = map[commands.Command]CommandHandler{
		commands.Pick:     PickHandler,
		commands.Goto:     GotoHandler,
		commands.Lookup:   LookupHandler,
		commands.Light:    LightHandler,
		commands.Describe: DescribeHandler,
		commands.Drop:     DropHandler,
		commands.Select:   SelectHandler,
		commands.Use:      UseHandler,
	}
}

//NewRoom - constrictor
func NewRoom(name string, desc string, init func(*Room), rooms []string, gs actor.StreamInterface) Room {
	a := area.NewArea(name, gs)
	room := new(Room)
	room.Area = a
	room.Actor.ProcessEvent = room.ProcessEvent
	room.NPCs = make(map[string]*actor.NPCInterface)
	room.Objects = make(map[string]actor.ObjectInterface)
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
		a.SendCompleterList(event.Sender, string(commands.Goto), a.ToRooms)
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
	tokens := strings.Split(event.Payload.(string), " ")
	command := strings.ToLower(tokens[0])
	_, ok := a.Streams[event.Sender]
	log.Println(fmt.Sprintf("%v: Recv command %s", a.Name, event))
	cmd := commands.Command(command)
	if ok == false && cmd != commands.Login && event.Sender != "cmd" {
		log.Println("Discard command " + command + " from " + event.Sender)
		return
	}
	handler, ok := Handlers[cmd]
	if ok {
		handler(a, event, tokens)
		return
	}
	switch cmd {
	case "_routes":
		a.SendCompleterList(event.Sender, string(commands.Goto), a.ToRooms)
	case "_items":
		a.SendCompleterListItems(event.Sender, string(commands.Pick), a.Items.GetItems())
	case "routes":
		a.SendEvent(event.Sender, events.SYSTEMMESSAGE, a.ToRooms)
	case commands.Unselect:
		p := *a.GetPlayer(event.Sender)
		p.SetSelected(nil)
	case commands.Status:
		a.SendEvent(event.Sender, events.STATUS, nil)
	default:
		if strings.HasPrefix(command, string(commands.Say)) {
			a.Broadcast(events.MESSAGE, event.Payload.(string)[1:len(event.Payload.(string))], event.Sender)
		} else {
			a.ForwardEvent("global", event)
		}
	}
}

//GetDesc -
func (a *Room) GetDesc() string {
	return a.Desc
}

//AddItem -
func (a *Room) AddItem(item actor.ItemInterface) {
	pos := actor.Index(a.State.Items, item.GetName())
	a.Items.AddItem(item.GetName(), item)
	if pos == -1 {
		a.State.Items = append(a.State.Items, item.GetName())
		a.UpdateState()
	}
}

//RemoveItem -
func (a *Room) RemoveItem(name string) {
	a.Items.RemoveItem(name)
	pos := actor.Index(a.State.Items, name)
	a.State.Items = append(a.State.Items[:pos], a.State.Items[pos+1:]...)
	a.UpdateState()
}

//GetItem -
func (a *Room) GetItem(name string) (actor.ItemInterface, bool) {
	return a.Items.GetItem(name)
}

func (a *Room) Inspect() string {
	return fmt.Sprintf("%v\n%v", a.Desc, a.Items)
}
