package time_stream

import (
	"actor"
	"events"
	"time"
	// "fmt"
)

type Stream struct {
	actor.Actor
}

func NewStream(gs chan events.Event) *Stream {
	a := actor.NewActor(gs)
	actor := new(Stream)
	actor.Actor = *a
	return actor
}

func (a Stream) Live() {
	k := 1
	ticks := 0
	ticker := time.NewTicker(time.Duration(100 * k * int(time.Millisecond)))
	for t := range ticker.C {
		a.Streams["global"] <- events.Event{t, events.TICK, nil}
		// log.Println(ticks, ticks%10)
		if ticks%10 == 0 {
			a.Streams["global"] <- events.Event{t, events.SECOND, nil}
		}
		if ticks%600 == 0 {
			a.Streams["global"] <- events.Event{t, events.MINUTE, nil}
		}
		ticks++
	}
}
