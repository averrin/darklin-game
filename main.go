package main

import (
	core "core"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"

	"gopkg.in/mgo.v2"
)

func main() {
	var err error
	core.SESSION, err = mgo.Dial("mongo")
	if err != nil {
		panic(err)
	}
	defer core.SESSION.Close()
	core.SESSION.SetMode(mgo.Monotonic, true)

	gs := core.NewGlobalStream()
	core.WORLD = core.NewWorld(&gs)
	core.WORLD.Init()

	announcer := core.NewAnnouncer(&gs.Stream)
	go announcer.Live()

	// gs.Subscribe(SECOND, announcer)
	gs.Subscribe(core.MINUTE, &announcer.Actor)

	gs.Streams["time"] = &core.WORLD.Time.Stream

	http.HandleFunc("/ws", gs.GetPlayerHandler())

	port := flag.Int("port", 80, "port for serving")
	interactive := flag.Bool("interactive", false, "readline mode")
	// debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()
	log.Println(fmt.Sprintf("Serving at :%v", *port))
	// http.Handle("/", http.FileServer(http.Dir(".")))
	go http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
	if *interactive == false {
		log.SetOutput(ioutil.Discard)
		gs.Live()
	} else {
		go gs.Live()
		core.RunShell(&gs.Stream)
	}
}
