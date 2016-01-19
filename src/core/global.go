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
	"github.com/fatih/color"

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

type Formatter struct {
	Blue   func(...interface{}) string
	Yellow func(...interface{}) string
	Red    func(...interface{}) string
	Green  func(...interface{}) string
}

// GlobalStream for global events
type GlobalStream struct {
	Actor
	Players   map[*websocket.Conn]*Player
	Storage   *Storage
	State     GlobalState
	Formatter Formatter
}

func (a *GlobalStream) GetPlayer(name string) *Player {
	for _, v := range a.Players {
		if v.Name == name {
			return v
		}
	}
	return &Player{}
}

// NewGlobalStream constructor
func NewGlobalStream() *GlobalStream {
	gs := make(chan Event)
	a := NewActor("global", gs)
	actor := new(GlobalStream)
	actor.Actor = *a
	actor.Players = make(map[*websocket.Conn]*Player)
	actor.Storage = NewStorage()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	formatter := Formatter{blue, yellow, red, green}
	actor.Formatter = formatter
	s := actor.Storage.Session.Copy()
	defer s.Close()
	db := s.DB("darklin")
	n, _ := db.C("state").Count()
	actor.State = *new(GlobalState)
	actor.State.Date = time.Date(774, 1, 1, 12, 0, 0, 0, time.UTC)
	actor.State.New = true
	if n != 0 {
		db.C("state").Find(bson.M{}).One(&actor.State)
		actor.State.New = false
	}
	return actor
}

func (a *GlobalStream) ProcessEvent(event Event) {
	formatter := a.Formatter
	blue := formatter.Blue
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
		a.Broadcast(MESSAGE, event.Payload, event.Sender)
	case COMMAND:
		log.Println(fmt.Sprintf("%v > %v", blue(event.Sender), event.Payload))
		a.ProcessCommand(event)
	}
}

func (a *GlobalStream) ProcessCommand(event Event) {
	formatter := a.Formatter
	blue := formatter.Blue
	tokens := strings.Split(event.Payload.(string), " ")
	log.Println(tokens, len(tokens))
	command := strings.ToLower(tokens[0])
	_, ok := a.Streams[event.Sender]
	if ok == false && command != "login" {
		return
	}
	switch command {
	case "reset":
		if event.Sender == "cmd" {
			a.SendEvent("time", RESET, a.Streams[event.Sender])
		}
	case "time":
		a.SendEvent("time", INFO, a.Streams[event.Sender])
	case "online":
		log.Println(fmt.Sprintf("Online: %v", len(a.Players)))
		if event.Sender != "cmd" {
			a.SendEvent(event.Sender, MESSAGE, fmt.Sprintf("Online: %v", len(a.Players)))
		}
	case "exit":
		os.Exit(0)
	case "login":
		if len(tokens) == 3 {
			log.Println("try login", blue(tokens[1]), tokens[2])
			_, ok := a.Streams[tokens[1]]
			if ok {
				var player *Player
				for _, p := range a.Players {
					if p.Name == event.Sender {
						player = p
						break
					}
				}
				player.Message(NewEvent(MESSAGE, "Пользователь с таким именем уже залогинен", "global"))
			} else {
				p := a.GetPlayer(event.Sender)
				// delete(a.Streams, p.Name)
				p.Name = tokens[1]
				a.Streams[p.Name] = p.Stream
				p.Loggedin = true
				go p.Live()
				a.SendEvent(p.Name, MESSAGE, "Вы вошли как: "+p.Name)
			}
		}
	default:
		if strings.HasPrefix(command, "/") {
			a.Broadcast(MESSAGE, event.Payload.(string)[1:len(event.Payload.(string))], event.Sender)
		}
	}
}

// Live method for dispatch events
func (a GlobalStream) Live() {
	s := a.Storage.Session.Copy()
	defer s.Close()
	a.Storage.DB = s.DB("darklin")
	for {
		event := <-a.Stream
		// log.Println(event)
		a.NotifySubscribers(event)
		a.ProcessEvent(event)
	}
}

// CmdHandler - handle user input
func (a GlobalStream) CmdHandler(w http.ResponseWriter, r *http.Request) {
	formatter := a.Formatter
	red := formatter.Red
	c, err := upgrader.Upgrade(w, r, nil)
	name := uuid.New()
	p := NewPlayer(name, a.Stream)
	p.Connection = c
	p.Message(NewEvent(MESSAGE, "Подключено. Наберите: login <username> <password>", "global"))
	a.Players[c] = p
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
			p.Loggedin = false
			p.Stream <- Event{time.Now(), CLOSE, nil, a.Name}
			delete(a.Players, c)
			delete(a.Streams, p.Name)
			break
		}
		line := string(message)
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		log.Printf("recv: %s", line)
		a.Stream <- Event{time.Now(), COMMAND, line, p.Name}
		// err = c.WriteMessage(mt, []byte("U r "+name))
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}
}
