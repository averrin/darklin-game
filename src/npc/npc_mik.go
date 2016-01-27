package npc

import "time"
import "events"

//NewMik - nobody likes darkness
func NewMik(gs *chan *events.Event) NPC {
	mik := NewNPC("Mik Rori", gs, "Hall")
	// mik.ProcessEvent = mik.Mik

	mik.Handlers[events.ROOMCHANGED] = mik.MikRoomChanged
	mik.Handlers[events.ROOMENTER] = mik.MikRoomEnter
	mik.Handlers[MIK_SMOKE] = mik.MikSmoke
	mik.Handlers[MIK_CHANGEROOM] = mik.MikChangeRoom
	mik.Handlers[events.LIGHT] = mik.MikLight

	che := NewEvent(MIK_CHANGEROOM, nil, mik.Name)
	che.ID = "Mik_change_room"
	che.Every = 2 * time.Minute
	cme := NewEvent(MIK_SMOKE, nil, mik.Name)
	cme.Every = 10 * time.Minute
	go func() {
		mik.Stream <- che
		mik.Stream <- cme
	}()
	return mik
}

const (
	MIK_CHANGEROOM events.EventType = iota
	MIK_SMOKE
)

func (a *NPC) MikRoomChanged(event *events.Event) bool {
	if !a.Room.State.Light {
		a.BroadcastRoom(ceventsore.MESSAGE, "И тут темень!", a.Name, a.Room)
	}
	return false
}

func (a *NPC) MikRoomEnter(event *events.Event) bool {
	a.Room.SendEventWithSender(event.Sender, events.MESSAGE, "Привет.", a.Name)
	return false
}

func (a *NPC) MikSmoke(event *events.Event) bool {
	a.BroadcastRoom(events.SYSTEMMESSAGE, "*Мик закуривает трубку*", a.Name, a.Room)
	return false
}

func (a *NPC) MikChangeRoom(event *events.Event) bool {
	if a.State.Room == "Hall" {
		a.ChangeRoom(world.WORLD.Rooms["second"])
	} else {
		a.ChangeRoom(world.WORLD.Rooms["Hall"])
	}
	return false
}

func (a *NPC) MikLight(event *events.Event) bool {
	if !event.Payload.(bool) {
		a.BroadcastRoom(events.MESSAGE, "Эй, кто выключил свет?", a.Name, a.Room)
		a.BroadcastRoom(events.SYSTEMMESSAGE, "*шорох, шаги, чирканье спичек*", a.Name, a.Room)
		ne := NewEvent(events.COMMAND, "light on", a.Name)
		ne.ID = "Mik_light_on"
		ne.Delay = 5 * time.Second
		a.Room.Stream <- ne
	} else {
		ev, ok := a.Room.PendingEvents["Mik_light_on"]
		if ok {
			a.BroadcastRoom(events.MESSAGE, "То-то же!", a.Name, a.Room)
			ev.Abort = true
		} else {
			a.BroadcastRoom(events.MESSAGE, "Так лучше!", a.Name, a.Room)
		}
	}
	return false
}
