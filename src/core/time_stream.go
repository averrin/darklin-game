package main

import (
	"fmt"
	"time"
	// "fmt"
)

// TimeStream - ticker
type TimeStream struct {
	Actor
}

// NewTimeStream constructor
func NewTimeStream(gs chan Event) *TimeStream {
	a := NewActor("time", gs)
	actor := new(TimeStream)
	actor.Actor = *a
	return actor
}

// Live method
func (a TimeStream) Live() {
	k := 1
	ticks := 0
	ticker := time.NewTicker(time.Duration(100 * k * int(time.Millisecond)))
	for t := range ticker.C {
		go func() {
			event := <-a.Stream
			if event.Type == INFO {
				event.Payload.(chan Event) <- Event{t, MESSAGE, fmt.Sprintf("Ticks since start: %v", ticks), "time"}
			}
		}()
		a.SendEvent("global", TICK, nil)
		// log.Println(ticks, ticks%10)
		if ticks > 0 && ticks%10 == 0 {
			a.SendEvent("global", SECOND, nil)
		}
		if ticks > 0 && ticks%600 == 0 {
			a.SendEvent("global", MINUTE, nil)
		}
		ticks++
	}
}
