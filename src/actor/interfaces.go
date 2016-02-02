package actor

import (
	"events"
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
	GetPlayerHandler() func(http.ResponseWriter, *http.Request)
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
	GetName() string
	GetDesc() string
	RemoveNPC(string)
	AddNPC(NPCInterface)
	AddPlayer(PlayerInterface)
	RemovePlayer(PlayerInterface)
	SendEventWithSender(string, events.EventType, interface{}, string)
	GetPendingEvent(string) (*events.Event, bool)
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
}

//NPCInterface -
type NPCInterface interface {
	Live()
	GetName() string
	ChangeRoom(*RoomInterface)
	GetStream() *chan *events.Event
	SetStream(string, *chan *events.Event)
	SetRoom(RoomInterface)
	// AddItem(ItemInterface)
}

//ItemInterface -
type ItemInterface interface {
	GetName() string
	GetDesc() string
}

//ItemContainerInterface -
type ItemContainerInterface interface {
	GetItems() map[string]ItemInterface
	GetItem(string) (ItemInterface, bool)
	AddItem(string, ItemInterface)
	RemoveItem(string)
	Count() int
}
