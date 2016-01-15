package global

import (
	"actor"
	"events"
	"fmt"
	"log"
	"net/http"
	"os"
	"player"
	"time"

	"code.google.com/p/go-uuid/uuid"

	// "golang.org/x/net/websocket"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// Stream for global events
type Stream struct {
	actor.Actor
	Players map[*websocket.Conn]player.Player
}

// NewStream constructor
func NewStream() *Stream {
	gs := make(chan events.Event)
	a := actor.NewActor("global", gs)
	actor := new(Stream)
	actor.Actor = *a
	actor.Players = make(map[*websocket.Conn]player.Player)
	return actor
}

// Live method for dispatch events
func (a Stream) Live() {
	for {
		event := <-a.Stream
		// log.Println(event)
		a.NotifySubscribers(event)
		switch event.Type {
		case events.MESSAGE:
			a.ForwardEvent("player", event)
		// 	log.Println("MESSAGE: ", event.Payload)
		case events.COMMAND:
			log.Println("> ", event.Payload)
			switch event.Payload {
			case "time":
				a.SendEvent("time", events.INFO, a.Streams[event.Sender])
			case "online":
				a.SendEvent(event.Sender, events.MESSAGE, fmt.Sprintf("Online: %v", len(a.Players)))
			case "exit":
				os.Exit(0)
			default:
				a.Broadcast(events.MESSAGE, event.Payload, event.Sender)
			}
		}
	}
}

// CmdHandler - handle user input
func (a Stream) CmdHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	name := uuid.New()
	p := player.NewPlayer(name, a.Stream)
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
			log.Println("Disconnect", name)
			delete(a.Players, c)
			delete(a.Streams, name)
			break
		}
		log.Printf("recv: %s", message)
		a.Stream <- events.Event{time.Now(), events.COMMAND, string(message), name}
		// err = c.WriteMessage(mt, []byte("U r "+name))
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}
}
