package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
func (a TestActor) ConsumeEvent(event *Event) {
	a.Stream <- event
}

// NewTestActor because i, sucj in golang yet
func NewTestActor(gs *chan *Event) *TestActor {
	a := NewActor("announcer", gs)
	actor := new(TestActor)
	actor.Actor = *a
	return actor
}

// ProcessEvent - i need print something
func (a TestActor) ProcessEvent(event *Event) {
	switch event.Type {
	case SECOND:
		a.SendEvent("global", MESSAGE, "Every second, mister")
	case MINUTE:
		a.SendEvent("global", MESSAGE, "Every minute, boss")
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
	log.Println(gs.Stream)
	ts := NewTimeStream(&gs.Stream, gs.State.Date)
	go ts.Live()

	testActor := NewTestActor(&gs.Stream)
	go testActor.Live()

	// gs.Subscribe(SECOND, testActor)
	gs.Subscribe(MINUTE, testActor)

	gs.Streams["time"] = &ts.Stream

	http.HandleFunc("/ws", gs.GetPlayerHandler())

	port := flag.Int("port", 80, "port for serving")
	interactive := flag.Bool("interactive", false, "readline mode")
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()
	log.Println(fmt.Sprintf("Serving at :%v", *port))
	// http.Handle("/", http.FileServer(http.Dir(".")))
	go http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
	if *debug == true {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		log.Println(mem.Alloc)
		log.Println(mem.TotalAlloc)
		log.Println(mem.HeapAlloc)
		log.Println(mem.HeapSys)
	}
	if *interactive == false {
		log.SetOutput(ioutil.Discard)
		gs.Live()
	} else {
		go gs.Live()
		RunShell(&gs.Stream)
	}
}
