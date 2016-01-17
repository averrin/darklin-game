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
func NewTimeStream(gs chan Event, date time.Time) *TimeStream {
	a := NewActor("time", gs)
	actor := new(TimeStream)
	actor.Actor = *a
	actor.Date = date
	return actor
}

// Live method
func (a TimeStream) Live() {
	k := 1
	ticks := 0
	ticker := time.NewTicker(time.Duration(100 * k * int(time.Millisecond)))
	go func() {
		for {
			event := <-a.Stream
			if event.Type == INFO {
				event.Payload.(chan Event) <- Event{time.Now(), MESSAGE, fmt.Sprintf("Time: %v", a.Date), "time"}
			}
		}
	}()
	for _ := range ticker.C {
		a.Date = a.Date.Add(time.Duration(100 * k * int(time.Millisecond)))
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
