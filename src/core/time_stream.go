package core

import "time"

// "fmt"

// TimeStream - ticker
type TimeStream struct {
	Actor
	Date  time.Time
	Speed int
	Ticks int
}

// NewTimeStream constructor
func NewTimeStream(gs *chan *Event, date time.Time) *TimeStream {
	a := actor.NewActor("time", gs)
	actor := new(TimeStream)
	actor.Actor = *a
	actor.Date = date
	actor.Speed = 1
	return actor
}

// Live method
func (a *TimeStream) Live() {
	a.Ticks = 0
	ticker := time.NewTicker(time.Duration(100 * a.Speed * int(time.Millisecond)))
	paused := false
	go func() {
		for {
			event := <-a.Stream
			switch event.Type {
			case INFO:
				event.Payload.(chan *Event) <- NewEvent(SYSTEMMESSAGE, a.Date.Format("Mon Jan _2 15:04:05 2006"), "time")
			case RESET:
				a.Date = time.Date(774, 1, 1, 12, 0, 0, 0, time.UTC)
			case PAUSE:
				paused = true
			}
		}
	}()
	for _ = range ticker.C {
		if !paused {
			a.Date = a.Date.Add(time.Duration(100 * a.Speed * int(time.Millisecond)))
			a.SendEvent("global", TICK, a.Date)
			if a.Ticks > 0 && a.Ticks%10 == 0 {
				a.SendEvent("global", SECOND, a.Date)
			}
			if a.Ticks > 0 && a.Ticks%600 == 0 {
				a.SendEvent("global", MINUTE, a.Date)
			}
			if a.Ticks > 0 && a.Ticks%(60*600) == 0 {
				a.SendEvent("global", HOUR, a.Date)
			}
			if a.Ticks > 0 && a.Ticks%(60*600*24) == 0 {
				a.SendEvent("global", DAY, a.Date)
				a.Ticks = 0
			}
			a.Ticks++
		}
	}
}

func (a *TimeStream) Sleep(duration time.Duration) {
	start := a.Ticks
	tick := 100 * a.Speed * int(time.Millisecond)
	for time.Duration(a.Ticks*tick) < time.Duration(start*tick+int(duration)) {
		time.Sleep(time.Millisecond * 20)
	}
}

var TIME *TimeStream
