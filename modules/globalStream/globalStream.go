package globalStream

import (
	actor "../actor"
	area "../area"
	commands "../commands"
	events "../events"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/pborman/uuid"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// GlobalState -- info about world
type GlobalState struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Date time.Time
	New  bool
}

// GlobalStream for global events
type GlobalStream struct {
	area.Area
	State     GlobalState
	NewPlayer func(string, actor.StreamInterface) actor.PlayerInterface
}

// NewGlobalStream constructor
func NewGlobalStream() *GlobalStream {
	stream := new(GlobalStream)
	stream.Stream = make(chan *events.Event, 100)
	a := area.NewArea("global", stream)
	stream.Area = a
	s := stream.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	n, _ := db.C("state").Count()
	stream.State = *new(GlobalState)
	stream.State.Date = time.Date(774, 1, 1, 12, 0, 0, 0, time.UTC)
	stream.State.New = true
	stream.Actor.ProcessEvent = stream.ProcessEvent
	stream.Actor.ProcessCommand = stream.ProcessCommand
	if n != 0 {
		db.C("state").Find(bson.M{}).One(&stream.State)
		stream.State.New = false
		stream.State.Date = stream.State.Date.In(time.UTC)
	}
	return stream
}

//ProcessEvent in global stream
func (a *GlobalStream) ProcessEvent(event *events.Event) {
	formatter := a.Formatter
	// blue := formatter.Blue
	yellow := formatter.Yellow
	switch event.Type {
	case events.SECOND:
		if a.State.New {
			i := bson.NewObjectId()
			go a.Storage.DB.C("state").Insert(bson.M{"_id": i}, a.State)
			// log.Println(err)
			a.State.New = false
			a.State.ID = i
		}
		a.Broadcast(events.HEARTBEAT, event.Payload, "heartbeat")
		a.State.Date = event.Payload.(time.Time)
		go func() {
			_ = a.Storage.DB.C("state").Update(bson.M{"_id": a.State.ID}, a.State)
			// log.Println(err)
		}()
	case events.MESSAGE:
		log.Println(yellow("MESSAGE:"), event.Payload)
		if event.Sender != "Announcer" {
			p := *a.GetPlayer(event.Sender)
			room, _ := p.GetRoom()
			(*room).BroadcastRoom(events.MESSAGE, event.Payload, event.Sender)
		} else {
			a.Broadcast(events.MESSAGE, event.Payload, event.Sender)
		}
	case events.LOGGEDIN:
		p := *a.GetPlayer(event.Sender)
		a.Streams[p.GetName()] = p.GetStream()
	case events.COMMAND:
		// log.Println(fmt.Sprintf("%v > %v", blue(event.Sender), event.Payload))
		a.ProcessCommand(event)
	}
}

//ProcessCommand from user or cmd
func (a *GlobalStream) ProcessCommand(event *events.Event) {
	// log.Println(event)
	// formatter := a.Formatter
	// blue := formatter.Blue
	tokens := strings.Split(event.Payload.(string), " ")
	// log.Println(tokens, len(tokens))
	command := strings.ToLower(tokens[0])
	_, ok := a.Streams[event.Sender]
	// log.Println(ok, command, event)
	if ok == false && command != "login" && event.Sender != "cmd" {
		log.Println("Discard command " + command + " from " + event.Sender)
		return
	}
	switch commands.Command(command) {
	case "info":
		log.Println(fmt.Sprintf("Players: %v", a.Players))
		log.Println(fmt.Sprintf("Streams: %v", a.Streams))
		h, _ := a.World.GetRoom("Hall")
		log.Println(fmt.Sprintf("Hall: %v", *h))
		s, _ := a.World.GetRoom("Store")
		log.Println(fmt.Sprintf("Store: %v", *s))
	case "reset":
		if event.Sender == "cmd" {
			go a.SendEvent("time", events.RESET, a.Streams[event.Sender])
		}
	case "pause":
		if event.Sender == "cmd" {
			go a.SendEvent("time", events.PAUSE, nil)
		}
	case commands.Time:
		if event.Sender == "cmd" {
			panic("TODO: fix it")
			// log.Println(fmt.Sprintf("Date: %v", world.WORLD.Time.Date))
		} else {
			go a.SendEvent("time", events.INFO, *a.Streams[event.Sender])
		}
	case commands.Online:
		log.Println(fmt.Sprintf("Online: %v", len(a.Players)))
		if event.Sender != "cmd" {
			go a.SendEvent(event.Sender, events.SYSTEMMESSAGE, fmt.Sprintf("Online: %v", len(a.Players)))
		}
	case commands.Exit:
		os.Exit(0)
	case commands.Login:
		//TODO: do it faster
		// go func() {
		if len(tokens) == 3 {
			// log.Println("try login", blue(tokens[1]), tokens[2])
			_, ok := a.Streams[tokens[1]]
			if ok {
				player := *a.GetPlayer(event.Sender)
				player.Message(events.NewEvent(events.ERROR, "Пользователь с таким именем уже залогинен", "global"))
			} else {
				p := *a.GetPlayer(event.Sender)
				go p.ProcessEvent(events.NewEvent(events.LOGIN, tokens, "global"))
			}
		}
		// }()
	case commands.Help:
		p := *a.GetPlayer(event.Sender)
		go a.SendEvent(p.GetName(), events.SYSTEMMESSAGE, "Help message")
	default:
		if strings.HasPrefix(command, "/") {
			go a.Broadcast(events.MESSAGE, event.Payload.(string)[1:len(event.Payload.(string))], event.Sender)
		}
	}
}

// GetPlayerHandler - handle user input
func (a *GlobalStream) GetPlayerHandler(version string) func(w http.ResponseWriter, r *http.Request) {
	formatter := a.Formatter
	red := formatter.Red

	return func(w http.ResponseWriter, r *http.Request) {
		name := uuid.New()
		p := a.NewPlayer(name, a)
		p.SetStream("room", &a.Stream)
		c, err := upgrader.Upgrade(w, r, nil)
		p.SetConnection(c)
		p.Message(events.NewEvent(events.CONNECTED, version, "global"))
		p.Message(events.NewEvent(events.SYSTEMMESSAGE, "Подключено. Наберите: login <username> <password>", "global"))
		a.Players[p] = c
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				log.Println(red("Disconnect"), name)
				// p.Loggedin = false
				go func() {
					stream := *p.GetStream()
					stream <- events.NewEvent(events.CLOSE, nil, a.Name)
				}()
				room, _ := p.GetRoom()
				(*room).RemovePlayer(p)
				delete(a.Players, p)
				delete(a.Streams, p.GetName())
				return
			}
			line := string(message)
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			room, ok := p.GetRoom()
			if ok {
				stream := *(*room).GetStream()
				stream <- events.NewEvent(events.COMMAND, line, p.GetName())
			} else {
				a.Stream <- events.NewEvent(events.COMMAND, line, p.GetName())
			}
		}
	}
}

//GetDate -
func (a *GlobalStream) GetDate() time.Time {
	return a.State.Date
}
