package global

import (
	"actor"
	"events"
	"log"
)

// Stream for global events
type Stream struct {
	actor.Actor
}

func NewStream() *Stream {
	gs := make(chan events.Event)
	a := actor.NewActor(gs)
	actor := new(Stream)
	actor.Actor = *a
	return actor
}

// Live method for dispatch events
func (a Stream) Live() {
	for {
		event := <-a.Stream
		// log.Println(event)
		for _, s := range a.Subscriptions {
			if event.Type == s.Type || s.Type == events.ALL {
				go s.Subscriber.ConsumeEvent(event)
			}
		}
		switch event.Type {
		case events.MESSAGE:
			a.SendEvent("player", events.MESSAGE, event.Payload)
		// 	log.Println("MESSAGE: ", event.Payload)
		case events.COMMAND:
			log.Println("USER COMMAND: ", event.Payload)
			switch event.Payload {
			case "time":
				a.SendEvent("time", events.INFO, a.Streams["player"])
			}
		}
	}
}
