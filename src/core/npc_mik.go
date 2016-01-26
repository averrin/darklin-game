package main

import (
	"log"
	"time"
)

//NewMik - nobody likes darkness
func NewMik(gs *chan *Event) NPC {
	room := WORLD.Rooms["first"]
	mik := NewNPC("Mik Rori", gs, room)
	mik.ProcessEvent = mik.Mik
	che := NewEvent(MIK_CHANGEROOM, nil, mik.Name)
	che.ID = "Mik_change_room"
	che.Every = 1 * time.Minute
	go func() {
		mik.Stream <- che
	}()
	return mik
}

const (
	MIK_CHANGEROOM EventType = iota
)

//Mik - Mik event loop
func (a *NPC) Mik(event *Event) {
	log.Println(event)
	switch event.Type {
	case ROOMCHANGED:
		if !a.Room.State.Light {
			a.BroadcastRoom(MESSAGE, "И тут темень!", a.Name, a.Room)
		}
	case ROOMENTER:
		a.Room.SendEventWithSender(event.Sender, MESSAGE, "Привет.", a.Name)
	case MIK_CHANGEROOM:
		if a.State.Room == "first" {
			a.ChangeRoom(WORLD.Rooms["second"])
		} else {
			a.ChangeRoom(WORLD.Rooms["first"])
		}
	case LIGHT:
		if !event.Payload.(bool) {
			a.BroadcastRoom(MESSAGE, "Эй, кто выключил свет?", a.Name, a.Room)
			a.BroadcastRoom(SYSTEMMESSAGE, "*шорох, шаги, чирканье спичек*", a.Name, a.Room)
			ne := NewEvent(COMMAND, "light on", a.Name)
			ne.ID = "Mik_light_on"
			ne.Delay = 5 * time.Second
			a.Room.Stream <- ne
		} else {
			ev, ok := a.Room.PendingEvents["Mik_light_on"]
			if ok {
				a.BroadcastRoom(MESSAGE, "То-то же!", a.Name, a.Room)
				ev.Abort = true
			} else {
				a.BroadcastRoom(MESSAGE, "Так лучше!", a.Name, a.Room)
			}
		}
	}
}
