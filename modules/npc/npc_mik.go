package npc

import "time"

import events "../events"
import actor "../actor"

//NewMik - nobody likes darkness
func NewMik(gs actor.StreamInterface) *NPC {
	mik := NewNPC("Mik Rori", "Старик с трубкой.", gs)
	// mik.ProcessEvent = mik.Mik

	mik.Handlers[events.ROOMCHANGED] = mik.MikRoomChanged
	mik.Handlers[events.ROOMENTER] = mik.MikRoomEnter
	mik.Handlers[MIK_SMOKE] = mik.MikSmoke
	mik.Handlers[MIK_CHANGEROOM] = mik.MikChangeRoom
	mik.Handlers[events.LIGHT] = mik.MikLight

	che := events.NewEvent(MIK_CHANGEROOM, nil, mik.Name)
	che.ID = "Mik_change_room"
	che.Every = 1 * time.Minute
	// che.Every = 15 * time.Second
	cme := events.NewEvent(MIK_SMOKE, nil, mik.Name)
	cme.Every = 10 * time.Minute
	go func() {
		mik.Stream <- che
		mik.Stream <- cme
	}()
	return &mik
}

//Events
const (
	MIK_CHANGEROOM events.EventType = iota
	MIK_SMOKE
)

//MikRoomChanged -
func (a *NPC) MikRoomChanged(event *events.Event) bool {
	room := *a.Room
	if !room.GetState().Light {
		room.BroadcastRoom(events.MESSAGE, "И тут темень!", a.Name)
	}
	return false
}

//MikRoomEnter -
func (a *NPC) MikRoomEnter(event *events.Event) bool {
	room := *a.Room
	// log.Println(room)
	room.SendEventWithSender(event.Sender, events.MESSAGE, "Привет.", a.Name)
	return false
}

//MikSmoke -
func (a *NPC) MikSmoke(event *events.Event) bool {
	room := *a.Room
	room.BroadcastRoom(events.SYSTEMMESSAGE, "*Мик закуривает трубку*", a.Name)
	return false
}

//MikChangeRoom -
func (a *NPC) MikChangeRoom(event *events.Event) bool {
	world := a.World
	// log.Println(world, a.State.New)
	if a.State.Room == "Hall" {
		room, _ := world.GetRoom("Store")
		a.ChangeRoom(room)
	} else {
		room, _ := world.GetRoom("Hall")
		a.ChangeRoom(room)
	}
	return false
}

//MikLight -
func (a *NPC) MikLight(event *events.Event) bool {
	room := *a.Room
	if !event.Payload.(bool) {
		room.BroadcastRoom(events.MESSAGE, "Эй, кто выключил свет?", a.Name)
		room.BroadcastRoom(events.SYSTEMMESSAGE, "*шорох, шаги, чирканье спичек*", a.Name)
		ne := events.NewEvent(events.COMMAND, "light on", a.Name)
		ne.ID = "Mik_light_on"
		ne.Delay = 5 * time.Second
		stream := *room.GetStream()
		stream <- ne
	} else {
		ev, ok := room.GetPendingEvent("Mik_light_on")
		if ok {
			room.BroadcastRoom(events.MESSAGE, "То-то же!", a.Name)
			ev.Abort = true
		} else {
			room.BroadcastRoom(events.MESSAGE, "Так лучше!", a.Name)
		}
	}
	return false
}
