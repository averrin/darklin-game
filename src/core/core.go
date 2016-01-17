package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"gopkg.in/mgo.v2"
)

// TestActor just someone who do something
type TestActor struct {
	Actor
}

// ConsumeEvent of couse
func (a TestActor) ConsumeEvent(event Event) {
	a.Stream <- event
}

// NewTestActor because i, sucj in golang yet
func NewTestActor(gs chan Event) *TestActor {
	a := NewActor("announcer", gs)
	actor := new(TestActor)
	actor.Actor = *a
	return actor
}

// Live - i need print something
func (a TestActor) Live() {
	for {
		event := <-a.Stream
		a.NotifySubscribers(event)
		switch event.Type {
		case SECOND:
			a.SendEvent("global", MESSAGE, "Every second, mister")
		case MINUTE:
			a.SendEvent("global", MESSAGE, "Every minute, boss")
		}
	}
}

func main() {
	var err error
	session, err = mgo.Dial("mongo")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	gs := NewGlobalStream()
	stream := gs.Stream
	ts := NewTimeStream(stream, gs.State.Date)
	go ts.Live()

	testActor := NewTestActor(stream)
	go testActor.Live()

	// gs.Subscribe(SECOND, testActor)
	gs.Subscribe(MINUTE, testActor)

	gs.Streams["time"] = ts.Stream

	http.HandleFunc("/ws", gs.CmdHandler)

	port := flag.Int("port", 80, "port for serving")
	interactive := flag.Bool("interactive", false, "readline mode")
	flag.Parse()
	log.Println(fmt.Sprintf("Serving at :%v", *port))
	// http.Handle("/", http.FileServer(http.Dir(".")))
	go http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	if *interactive == false {
		gs.Live()
	} else {
		go gs.Live()
		RunShell(stream)
	}
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	log.Println(mem.Alloc)
	log.Println(mem.TotalAlloc)
	log.Println(mem.HeapAlloc)
	log.Println(mem.HeapSys)
}
