package global

import (
	"actor"
	"events"
	"fmt"
)

// Stream for global events
type Stream struct {
	actor.Actor
}


// Live method for dispatch events
func (a Stream) Live() {
	for {
		event := <- a.Stream
		for _, s := range a.Subscriptions {
			if event.Type == s.Type {
				go s.Subscriber.ConsumeEvent(event)
			}
		}
		switch event.Type {
		case events.MESSAGE:
			fmt.Println("MESSAGE: ", event.Payload)
		}
	}
}
