package global

import (
	"actor"
	"events"
	"fmt"
)

type GlobalStream struct {
	actor.Actor
}

func (actor GlobalStream) Live() {
	for {
		event := <-actor.Stream
		for _, s := range actor.Subscriptions {
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
