package actor

import (
	events "../events"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// EventPublisher - can send
type EventPublisher interface {
	SendEvent(events.EventType, interface{})
}

//StreamInterface -
type StreamInterface interface {
	Live()
	SetWorld(WorldInterface)
	GetWorld() WorldInterface
	GetStream() *chan *events.Event
	SetStream(string, *chan *events.Event)
	GetDate() time.Time
	Subscribe(events.EventType, *Actor)
	GetPlayerHandler(string) func(http.ResponseWriter, *http.Request)
}

//WorldInterface -
type WorldInterface interface {
	GetRoom(string) (*RoomInterface, bool)
	GetGlobal() *StreamInterface
	GetTime() *TimeInterface
	GetDate() time.Time
	AddRoom(string, RoomInterface)
	AddItem(ItemInterface)
	RemoveItem(string)
	GetItem(string) (ItemInterface, bool)
}

//TimeInterface -
type TimeInterface interface {
	Live()
	Sleep(time.Duration)
	GetDate() time.Time
	GetStream() *chan *events.Event
}

//RoomInterface -
type RoomInterface interface {
	Live()
	// Init()
	BroadcastRoom(events.EventType, interface{}, string)
	GetStream() *chan *events.Event
	GetState() AreaState
	RemoveNPC(string)
	AddNPC(NPCInterface)
	AddPlayer(PlayerInterface)
	RemovePlayer(PlayerInterface)
	SendEventWithSender(string, events.EventType, interface{}, string)
	GetPendingEvent(string) (*events.Event, bool)
	SelectableInterface
}

//PlayerInterface -
type PlayerInterface interface {
	Live()
	GetName() string
	GetStream() *chan *events.Event
	SetStream(string, *chan *events.Event)
	GetRoom() (*RoomInterface, bool)
	ChangeRoom(*RoomInterface)
	Message(*events.Event)
	ProcessEvent(*events.Event)
	SetConnection(*websocket.Conn)
	GetConnection() *websocket.Conn
	AddItem(ItemInterface)
	RemoveItem(string)
	GetItem(string) (ItemInterface, bool)
	GetItems() map[string]ItemInterface
	SetSelected(*SelectableInterface)
	GetSelected() *SelectableInterface
}

//NPCInterface -
type NPCInterface interface {
	Live()
	ChangeRoom(*RoomInterface)
	GetStream() *chan *events.Event
	SetStream(string, *chan *events.Event)
	SetRoom(RoomInterface)
	// AddItem(ItemInterface)
	SelectableInterface
}

//ItemInterface -
type ItemInterface interface {
	UsableInterface
	SelectableInterface
}

//ItemContainerInterface -
type ItemContainerInterface interface {
	GetItems() map[string]ItemInterface
	GetItem(string) (ItemInterface, bool)
	AddItem(string, ItemInterface)
	RemoveItem(string)
	Count() int
	String() string
}

//ObjectInterface -
type ObjectInterface interface {
	GetItems() map[string]ItemInterface
	GetItem(string) (ItemInterface, bool)
	AddItem(ItemInterface)
	RemoveItem(string)
	GetState() interface{}
	SelectableInterface
	UsableInterface
}

//SelectableInterface -
type SelectableInterface interface {
	GetName() string
	GetDesc() string
	Inspect() string
	// Select()
	// Unselect()
}

type UsableInterface interface {
	Use(ItemInterface) interface{}
}
