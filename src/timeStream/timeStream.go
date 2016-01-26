package timeStream

import "time"
import "events"
import "actor"

// "fmt"

// TimeStream - ticker
type TimeStream struct {
	actor.Actor
	Date  time.Time
	Speed int
	Ticks int
}

// NewTimeStream constructor
func NewTimeStream(gs *chan *events.Event, date time.Time) *TimeStream {
	a := actor.NewActor("time", gs)
	stream := new(TimeStream)
	stream.Actor = *a
	stream.Date = date
	stream.Speed = 1
	return stream
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
			case events.INFO:
				event.Payload.(chan *Event) <- NewEvent(events.SYSTEMMESSAGE, a.Date.Format("Mon Jan _2 15:04:05 2006"), "time")
			case events.RESET:
				a.Date = time.Date(774, 1, 1, 12, 0, 0, 0, time.UTC)
			case events.PAUSE:
				paused = true
			}
		}
	}()
	for _ = range ticker.C {
		if !paused {
			a.Date = a.Date.Add(time.Duration(100 * a.Speed * int(time.Millisecond)))
			a.SendEvent("global", events.TICK, a.Date)
			if a.Ticks > 0 && a.Ticks%10 == 0 {
				a.SendEvent("global", events.SECOND, a.Date)
			}
			if a.Ticks > 0 && a.Ticks%600 == 0 {
				a.SendEvent("global", events.MINUTE, a.Date)
			}
			if a.Ticks > 0 && a.Ticks%(60*600) == 0 {
				a.SendEvent("global", events.HOUR, a.Date)
			}
			if a.Ticks > 0 && a.Ticks%(60*600*24) == 0 {
				a.SendEvent("global", events.DAY, a.Date)
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
