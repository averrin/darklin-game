package rooms

import (
	"actor"
	"events"
	"fmt"
)

func PickHandler(a *Room, event *events.Event, tokens []string) {
	if len(tokens) == 2 {
		p := *a.GetPlayer(event.Sender)
		if p.GetSelected() == nil {
			item, ok := a.Items.GetItem(tokens[1])
			if ok {
				a.RemoveItem(tokens[1])
				p.AddItem(item)
				go a.SendEvent(event.Sender, events.SYSTEMMESSAGE, fmt.Sprintf("Вы подняли: %v [%v]", item.GetDesc(), item.GetName()))
				go a.SendCompleterListItems(event.Sender, "drop", p.GetItems())
				go a.SendCompleterListItems(event.Sender, "pick", a.Items.GetItems())
			}
		} else {
			obj := p.GetSelected()
			switch (*obj).(type) {
			case actor.ObjectInterface:
				sel := (*obj).(actor.ObjectInterface)
				item, ok := sel.GetItem(tokens[1])
				if ok {
					sel.RemoveItem(tokens[1])
					p.AddItem(item)
					go a.SendEvent(event.Sender, events.SYSTEMMESSAGE, fmt.Sprintf("Вы подняли: %v [%v]", item.GetDesc(), item.GetName()))
					go a.SendCompleterListItems(event.Sender, "drop", p.GetItems())
					go a.SendCompleterListItems(event.Sender, "pick", a.Items.GetItems())
				}
			}
		}
		go a.UpdateState()
	}
}

func GotoHandler(a *Room, event *events.Event, tokens []string) {
	go func() {
		if len(tokens) == 2 {
			p := *a.GetPlayer(event.Sender)
			w := a.World
			room, ok := w.GetRoom(tokens[1])
			if ok && stringInSlice(tokens[1], a.ToRooms) {
				p.ChangeRoom(room)
			} else {
				a.SendEvent(event.Sender, events.ERROR, fmt.Sprintf("Вы не можете перейти в эту комнату: %v", tokens[1]))
			}
		}
	}()
}

func LookupHandler(a *Room, event *events.Event, tokens []string) {
	if a.State.Light {
		p := *a.GetPlayer(event.Sender)
		if p.GetSelected() == nil {
			a.SendEvent(event.Sender, events.DESCRIBE, a.Inspect())
			// a.SendEvent(event.Sender, events.DESCRIBE, fmt.Sprintf("Предметы: \n%v", a.Items))
			// a.SendEvent(event.Sender, events.DESCRIBE, fmt.Sprintf("Объекты: \n%v", a.Objects))
			// go a.SendCompleterListItems(event.Sender, "pick", a.Items.GetItems())
		} else {
			obj := p.GetSelected()
			switch (*obj).(type) {
			case actor.ObjectInterface:
				o := (*obj).(actor.ObjectInterface)
				a.SendEvent(event.Sender, events.DESCRIBE, o.Inspect())
				// a.SendEvent(event.Sender, events.DESCRIBE, fmt.Sprintf("Предметы: \n%v", o.GetItems()))
			}
		}
	} else {
		go a.SendEvent(event.Sender, events.SYSTEMMESSAGE, "В комнате темно")
	}
}

func LightHandler(a *Room, event *events.Event, tokens []string) {
	if len(tokens) == 2 && (tokens[1] == "on" || tokens[1] == "off") {
		if tokens[1] == "on" {
			if a.State.Light {
				go a.SendEvent(event.Sender, events.SYSTEMMESSAGE, "В комнате уже светло")
				return
			}
			a.State.Light = true
			go a.Broadcast(events.SYSTEMMESSAGE, "В комнате зажегся свет", a.Name)
		} else {
			if !a.State.Light {
				go a.SendEvent(event.Sender, events.SYSTEMMESSAGE, "В комнате уже темно")
				return
			}
			a.State.Light = false
			go a.Broadcast(events.SYSTEMMESSAGE, "В комнате погас свет", a.Name)
		}
		go func() { a.Stream <- events.NewEvent(events.LIGHT, a.State.Light, event.Sender) }()
		go a.Broadcast(events.LIGHT, a.State.Light, a.Name)
		go a.UpdateState()
	}
}

func DescribeHandler(a *Room, event *events.Event, tokens []string) {
	p := *a.GetPlayer(event.Sender)
	if p.GetSelected() == nil {
		a.SendEvent(event.Sender, events.DESCRIBE, a.Desc)
	} else {
		a.SendEvent(event.Sender, events.DESCRIBE, (*p.GetSelected()).GetDesc())
	}
}

func DropHandler(a *Room, event *events.Event, tokens []string) {
	if len(tokens) == 2 {
		p := *a.GetPlayer(event.Sender)
		item, ok := p.GetItem(tokens[1])
		if ok {
			a.AddItem(item)
			p.RemoveItem(tokens[1])
			go a.SendEvent(event.Sender, events.SYSTEMMESSAGE, fmt.Sprintf("Вы бросили: %v [%v]", item.GetDesc(), item.GetName()))
			go a.SendCompleterListItems(event.Sender, "drop", p.GetItems())
			go a.SendCompleterListItems(event.Sender, "pick", a.Items.GetItems())
		}
	}
}

func SelectHandler(a *Room, event *events.Event, tokens []string) {
	object, ok := a.Objects[tokens[1]]
	p := *a.GetPlayer(event.Sender)
	if ok {
		obj := object.(actor.SelectableInterface)
		p.SetSelected(&obj)
	}
}

func UseHandler(a *Room, event *events.Event, tokens []string) {
	if len(tokens) == 2 {
		p := *a.GetPlayer(event.Sender)
		item, ok := p.GetItem(tokens[1])
		if ok {
			obj := *p.GetSelected()
			if obj == nil {
				go a.SendEvent(event.Sender, events.DESCRIBE, item.Use(item))
			} else {
				switch obj.(type) {
				case actor.UsableInterface:
					o := obj.(actor.UsableInterface)
					a.SendEvent(event.Sender, events.DESCRIBE, o.Use(item))
				}
			}
			go a.UpdateState()
		}
	}
}
