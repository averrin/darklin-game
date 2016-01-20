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
func NewTimeStream(gs *chan *Event, date time.Time) *TimeStream {
	a := NewActor("time", gs)
	actor := new(TimeStream)
	actor.Actor = *a
	actor.Date = date
	return actor
}

// Live method
func (a *TimeStream) Live() {
	k := 1
	ticks := 0
	ticker := time.NewTicker(time.Duration(100 * k * int(time.Millisecond)))
	paused := false
	go func() {
		for {
			event := <-*a.Stream
			switch event.Type {
			case INFO:
				event.Payload.(chan *Event) <- NewEvent(MESSAGE, fmt.Sprintf("Time: %v", a.Date), "time")
			case RESET:
				a.Date = time.Date(774, 1, 1, 12, 0, 0, 0, time.UTC)
			case PAUSE:
				paused = true
			}
		}
	}()
	for _ = range ticker.C {
		if !paused {
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
}
