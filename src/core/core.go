package main

import (
	"expvar"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"gopkg.in/mgo.v2"
)

var (
	exp_events_processed = expvar.NewInt("events_processed")
)

func main() {
	var err error
	session, err = mgo.Dial("mongo")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	gs := NewGlobalStream()
	// log.Println(gs.Stream)
	ts := NewTimeStream(&gs.Stream, gs.State.Date)
	go ts.Live()

	announcer := NewAnnouncer(&gs.Stream)
	go announcer.Live()

	// gs.Subscribe(SECOND, announcer)
	gs.Subscribe(MINUTE, announcer)

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
