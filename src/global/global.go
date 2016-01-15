package global

import (
	"actor"
	"events"
	"log"
	"os"
)

// Stream for global events
type Stream struct {
	actor.Actor
}

func NewStream() *Stream {
	gs := make(chan events.Event)
	a := actor.NewActor("global", gs)
	actor := new(Stream)
	actor.Actor = *a
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
			log.Println("USER COMMAND: ", event.Payload)
			switch event.Payload {
			case "time":
				a.SendEvent("time", events.INFO, a.Streams["player"])
			case "exit":
				os.Exit(0)
			}
		}
	}
}
