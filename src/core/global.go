package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"code.google.com/p/go-uuid/uuid"

	// "golang.org/x/net/websocket"
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
	Area
	State GlobalState
}

//GetPlayer by name
func (a *GlobalStream) GetPlayer(name string) *Player {
	for v := range a.Players {
		if v.Name == name {
			return v
		}
	}
	return &Player{}
}

// NewGlobalStream constructor
func NewGlobalStream() GlobalStream {
	gs := make(chan *Event, 100)
	a := NewArea("global", &gs)
	actor := new(GlobalStream)
	actor.Area = a
	s := actor.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	n, _ := db.C("state").Count()
	actor.State = *new(GlobalState)
	actor.State.Date = time.Date(774, 1, 1, 12, 0, 0, 0, time.UTC)
	actor.State.New = true
	actor.Actor.ProcessEvent = actor.ProcessEvent
	actor.Actor.ProcessCommand = actor.ProcessCommand
	if n != 0 {
		db.C("state").Find(bson.M{}).One(&actor.State)
		actor.State.New = false
		actor.State.Date = actor.State.Date.In(time.UTC)
	}
	return *actor
}

//ProcessEvent in global stream
func (a *GlobalStream) ProcessEvent(event *Event) {
	formatter := a.Formatter
	// blue := formatter.Blue
	yellow := formatter.Yellow
	switch event.Type {
	case SECOND:
		if a.State.New {
			i := bson.NewObjectId()
			go a.Storage.DB.C("state").Insert(bson.M{"_id": i}, a.State)
			// log.Println(err)
			a.State.New = false
			a.State.ID = i
		}
		a.Broadcast(HEARTBEAT, event.Payload, "heartbeat")
		a.State.Date = event.Payload.(time.Time)
		go func() {
			_ = a.Storage.DB.C("state").Update(bson.M{"_id": a.State.ID}, a.State)
			// log.Println(err)
		}()
	case MESSAGE:
		log.Println(yellow("MESSAGE:"), event.Payload)
		if event.Sender != "Announcer" {
			p := a.GetPlayer(event.Sender)
			a.BroadcastRoom(MESSAGE, event.Payload, event.Sender, p.Room)
		} else {
			a.Broadcast(MESSAGE, event.Payload, event.Sender)
		}
	case LOGGEDIN:
		p := a.GetPlayer(event.Sender)
		a.Streams[p.Name] = &p.Stream
	case COMMAND:
		// log.Println(fmt.Sprintf("%v > %v", blue(event.Sender), event.Payload))
		a.ProcessCommand(event)
	}
}

//ProcessCommand from user or cmd
func (a *GlobalStream) ProcessCommand(event *Event) {
	// log.Println(event)
	// formatter := a.Formatter
	// blue := formatter.Blue
	tokens := strings.Split(event.Payload.(string), " ")
	// log.Println(tokens, len(tokens))
	command := strings.ToLower(tokens[0])
	_, ok := a.Streams[event.Sender]
	// log.Println("Recv command " + command + " from " + event.Sender)
	if ok == false && command != "login" && event.Sender != "cmd" {
		log.Println("Discard command " + command + " from " + event.Sender)
		return
	}
	switch command {
	case "info":
		log.Println(fmt.Sprintf("Players: %v", a.Players))
		log.Println(fmt.Sprintf("Streams: %v", a.Streams))
	case "reset":
		if event.Sender == "cmd" {
			go a.SendEvent("time", RESET, a.Streams[event.Sender])
		}
	case "pause":
		if event.Sender == "cmd" {
			go a.SendEvent("time", PAUSE, nil)
		}
	case "time":
		if event.Sender == "cmd" {
			log.Println(fmt.Sprintf("Date: %v", WORLD.Time.Date))
		} else {
			go a.SendEvent("time", INFO, *a.Streams[event.Sender])
		}
	case "online":
		log.Println(fmt.Sprintf("Online: %v", len(a.Players)))
		if event.Sender != "cmd" {
			go a.SendEvent(event.Sender, SYSTEMMESSAGE, fmt.Sprintf("Online: %v", len(a.Players)))
		}
	case "exit":
		os.Exit(0)
	case "goto":
		go func() {
			if len(tokens) == 2 {
				p := a.GetPlayer(event.Sender)
				room, ok := WORLD.Rooms[tokens[1]]
				if ok {
					if p.Room == room {
						a.SendEvent(event.Sender, ERROR, fmt.Sprintf("You are already here: %v", tokens[1]))
					} else {
						log.Println(p.Name, room)
						p.ChangeRoom(room)
					}
				} else {
					a.SendEvent(event.Sender, ERROR, fmt.Sprintf("No such room: %v", tokens[1]))
				}
			}
		}()
	case "login":
		//TODO: do it faster
		// go func() {
		if len(tokens) == 3 {
			// log.Println("try login", blue(tokens[1]), tokens[2])
			_, ok := a.Streams[tokens[1]]
			if ok {
				player := a.GetPlayer(event.Sender)
				player.Message(NewEvent(ERROR, "Пользователь с таким именем уже залогинен", "global"))
			} else {
				p := a.GetPlayer(event.Sender)
				go p.ProcessEvent(NewEvent(LOGIN, tokens, "global"))
			}
		}
		// }()
	case "help":
		p := a.GetPlayer(event.Sender)
		go a.SendEvent(p.Name, SYSTEMMESSAGE, "Help message")
	default:
		if strings.HasPrefix(command, "/") {
			go a.Broadcast(MESSAGE, event.Payload.(string)[1:len(event.Payload.(string))], event.Sender)
		}
	}
}

// GetPlayerHandler - handle user input
func (a *GlobalStream) GetPlayerHandler() func(w http.ResponseWriter, r *http.Request) {
	formatter := a.Formatter
	red := formatter.Red

	return func(w http.ResponseWriter, r *http.Request) {
		name := uuid.New()
		p := NewPlayer(name, &a.Stream)
		p.Streams["room"] = &a.Stream
		c, err := upgrader.Upgrade(w, r, nil)
		p.Connection = c
		p.Message(NewEvent(CONNECTED, nil, "global"))
		p.Message(NewEvent(SYSTEMMESSAGE, "Подключено. Наберите: login <username> <password>", "global"))
		a.Players[&p] = c
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
					p.Stream <- NewEvent(CLOSE, nil, a.Name)
				}()
				delete(a.Players, &p)
				delete(a.Streams, p.Name)
				return
			}
			line := string(message)
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			*p.Streams["room"] <- NewEvent(COMMAND, line, p.Name)
		}
	}
}
