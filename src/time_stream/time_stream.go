package time_stream

import (
	"actor"
	"events"
	"fmt"
	"time"
	// "fmt"
)

// Stream
type Stream struct {
	actor.Actor
}

func NewStream(gs chan events.Event) *Stream {
	a := actor.NewActor("time", gs)
	actor := new(Stream)
	actor.Actor = *a
	return actor
}

func (a Stream) Live() {
	k := 1
	ticks := 0
	ticker := time.NewTicker(time.Duration(100 * k * int(time.Millisecond)))
	for t := range ticker.C {
		go func() {
			event := <-a.Stream
			if event.Type == events.INFO {
				event.Payload.(chan events.Event) <- events.Event{t, events.MESSAGE, fmt.Sprintf("Ticks since start: %v", ticks), "time"}
			}
		}()
		a.SendEvent("global", events.TICK, nil)
		// log.Println(ticks, ticks%10)
		if ticks > 0 && ticks%10 == 0 {
			a.SendEvent("global", events.SECOND, nil)
		}
		if ticks > 0 && ticks%600 == 0 {
			a.SendEvent("global", events.MINUTE, nil)
		}
		ticks++
	}
}
