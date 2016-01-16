package main

import (
	"fmt"
	"time"
	// "fmt"
)

// TimeStream - ticker
type TimeStream struct {
	Actor
	Date time.Time
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
	a.Date = time.Date(774, 1, 1, 12, 0, 0, 0, time.Local)
	for t := range ticker.C {
		a.Date = a.Date.Add(time.Duration(100 * k * int(time.Millisecond)))
		go func() {
			event := <-a.Stream
			if event.Type == INFO {
				event.Payload.(chan Event) <- Event{t, MESSAGE, fmt.Sprintf("Time: %v", a.Date), "time"}
			}
		}()
		a.SendEvent("global", TICK, a.Date)
		// log.Println(ticks, ticks%10)
		if ticks > 0 && ticks%10 == 0 {
			a.SendEvent("global", SECOND, a.Date)
		}
		if ticks > 0 && ticks%600 == 0 {
			a.SendEvent("global", MINUTE, a.Date)
			ticks = 0
		}
		ticks++
	}
}
