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

// GlobalStream for global events
type GlobalStream struct {
	Actor
	Players map[*websocket.Conn]Player
	Storage *Storage
	State   GlobalState
}

// NewGlobalStream constructor
func NewGlobalStream() *GlobalStream {
	gs := make(chan Event)
	a := NewActor("global", gs)
	actor := new(GlobalStream)
	actor.Actor = *a
	actor.Players = make(map[*websocket.Conn]Player)
	actor.Storage = NewStorage()
	db := actor.Storage.Session.Copy().DB("darklin")
	n, _ := db.C("state").Count()
	actor.State = *new(GlobalState)
	actor.State.Date = time.Date(774, 1, 1, 12, 0, 0, 0, time.Local)
	actor.State.New = true
	if n != 0 {
		db.C("state").Find(bson.M{}).One(&actor.State)
		actor.State.New = false
	}
	return actor
}

// Live method for dispatch events
func (a GlobalStream) Live() {
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	db := a.Storage.Session.Copy().DB("darklin")
	for {
		event := <-a.Stream
		// log.Println(event)
		a.NotifySubscribers(event)
		switch event.Type {
		case SECOND:
			if a.State.New {
				i := bson.NewObjectId()
				go db.C("state").Insert(bson.M{"_id": i}, a.State)
				// log.Println(err)
				a.State.New = false
				a.State.ID = i
			}
			a.Broadcast(HEARTBEAT, event.Payload, "heartbeat")
			a.State.Date = event.Payload.(time.Time)
			go func() {
				_ = db.C("state").Update(bson.M{"_id": a.State.ID}, a.State)
				// log.Println(err)
			}()
		case MESSAGE:
			log.Println(yellow("MESSAGE:"), event.Payload)
			a.Broadcast(MESSAGE, event.Payload, event.Sender)
		case COMMAND:
			log.Println(fmt.Sprintf("%v > %v", blue(event.Sender), event.Payload))
			switch event.Payload {
			case "time":
				a.SendEvent("time", INFO, a.Streams[event.Sender])
			case "online":
				log.Println(fmt.Sprintf("Online: %v", len(a.Players)))
				if event.Sender != "cmd" {
					a.SendEvent(event.Sender, MESSAGE, fmt.Sprintf("Online: %v", len(a.Players)))
				}
			case "exit":
				os.Exit(0)
			default:
				switch event.Payload.(type) {
				case string:
					if strings.HasPrefix(event.Payload.(string), "/") {
						a.Broadcast(MESSAGE, event.Payload.(string)[1:len(event.Payload.(string))], event.Sender)
					}
				}
			}
		}
	}
}

// CmdHandler - handle user input
func (a GlobalStream) CmdHandler(w http.ResponseWriter, r *http.Request) {
	red := color.New(color.FgRed).SprintFunc()
	c, err := upgrader.Upgrade(w, r, nil)
	name := uuid.New()
	p := NewPlayer(name, a.Stream)
	p.Connection = c
	go p.Live()
	a.Players[c] = *p
	a.Streams[name] = p.Stream
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			// log.Println("read:", err)
			log.Println(red("Disconnect"), name)
			delete(a.Players, c)
			delete(a.Streams, name)
			break
		}
		line := string(message)
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		log.Printf("recv: %s", line)
		a.Stream <- Event{time.Now(), COMMAND, line, name}
		// err = c.WriteMessage(mt, []byte("U r "+name))
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}
}
