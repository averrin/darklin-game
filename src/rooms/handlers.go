package rooms

import (
	"actor"
	"events"
	"fmt"
	"strings"
)

//PickHandler -
func PickHandler(a *Room, event *events.Event, tokens []string) {
	if len(tokens) >= 2 {
		p := *a.GetPlayer(event.Sender)
		itemName := strings.Join(tokens[1:], " ")
		if p.GetSelected() == nil {
			item, ok := a.Items.GetItem(itemName)
			if ok {
				a.RemoveItem(itemName)
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
				item, ok := sel.GetItem(itemName)
				if ok {
					sel.RemoveItem(itemName)
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

//GotoHandler -
func GotoHandler(a *Room, event *events.Event, tokens []string) {
	go func() {
		if len(tokens) >= 2 {
			itemName := strings.Join(tokens[1:], " ")
			p := *a.GetPlayer(event.Sender)
			w := a.World
			room, ok := w.GetRoom(itemName)
			if ok && stringInSlice(itemName, a.ToRooms) {
				p.ChangeRoom(room)
			} else {
				a.SendEvent(event.Sender, events.ERROR, fmt.Sprintf("Вы не можете перейти в эту комнату: %v", itemName))
			}
		}
	}()
}

//LookupHandler -
func LookupHandler(a *Room, event *events.Event, tokens []string) {
	if a.State.Light {
		p := *a.GetPlayer(event.Sender)
		if p.GetSelected() == nil {
			a.SendEvent(event.Sender, events.DESCRIBE, a.Inspect())
			// a.SendEvent(event.Sender, events.DESCRIBE, fmt.Sprintf("Предметы: \n%v", a.Items))
			// a.SendEvent(event.Sender, events.DESCRIBE, fmt.Sprintf("Объекты: \n%v", a.Objects))
			go a.SendCompleterListItems(event.Sender, "pick", a.Items.GetItems())
			go a.SendCompleterListObjects(event.Sender, "select", a.GetObjects())
		} else {
			obj := p.GetSelected()
			switch (*obj).(type) {
			case actor.ObjectInterface:
				o := (*obj).(actor.ObjectInterface)
				a.SendEvent(event.Sender, events.DESCRIBE, o.Inspect())
				go a.SendCompleterListItems(event.Sender, "pick", o.GetItems())
				// a.SendEvent(event.Sender, events.DESCRIBE, fmt.Sprintf("Предметы: \n%v", o.GetItems()))
			}
		}
	} else {
		go a.SendEvent(event.Sender, events.SYSTEMMESSAGE, "В комнате темно")
	}
}

//LightHandler -
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

//DescribeHandler -
func DescribeHandler(a *Room, event *events.Event, tokens []string) {
	p := *a.GetPlayer(event.Sender)
	if p.GetSelected() == nil {
		a.SendEvent(event.Sender, events.DESCRIBE, a.Desc)
	} else {
		a.SendEvent(event.Sender, events.DESCRIBE, (*p.GetSelected()).GetDesc())
	}
}

//DropHandler -
func DropHandler(a *Room, event *events.Event, tokens []string) {
	if len(tokens) >= 2 {
		itemName := strings.Join(tokens[1:], " ")
		p := *a.GetPlayer(event.Sender)
		item, ok := p.GetItem(itemName)
		if ok {
			a.AddItem(item)
			p.RemoveItem(itemName)
			go a.SendEvent(event.Sender, events.SYSTEMMESSAGE, fmt.Sprintf("Вы бросили: %v [%v]", item.GetDesc(), item.GetName()))
			go a.SendCompleterListItems(event.Sender, "drop", p.GetItems())
			go a.SendCompleterListItems(event.Sender, "pick", a.Items.GetItems())
		}
	}
}

//SelectHandler -
func SelectHandler(a *Room, event *events.Event, tokens []string) {
	itemName := strings.Join(tokens[1:], " ")
	object, ok := a.Objects[itemName]
	p := *a.GetPlayer(event.Sender)
	if ok {
		obj := object.(actor.SelectableInterface)
		p.SetSelected(&obj)
	}
}

//UseHandler -
func UseHandler(a *Room, event *events.Event, tokens []string) {
	if len(tokens) >= 2 {
		itemName := strings.Join(tokens[1:], " ")
		p := *a.GetPlayer(event.Sender)
		item, ok := p.GetItem(itemName)
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
