package time_stream

import (
	"actor"
	"events"
	"time"
	// "fmt"
)

type TimeStream struct {
	actor.Actor
}

func (a TimeStream) Live() {
	ticker := time.NewTicker(time.Millisecond * 500)
	for t := range ticker.C {
		a.GlobalStream <- events.Event{t, events.TICK, nil}
	}
}
