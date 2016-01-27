package main

import (
	core "core"
	"flag"
	"fmt"
	"globalStream"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"player"
	"world"

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

	gs := globalStream.NewGlobalStream()
	gs.NewPlayer = player.NewPlayer
	world := world.NewWorld(gs)
	world.Init()

	gs.SetStream("time", (*world.GetTime()).GetStream())

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
		core.RunShell(gs.GetStream())
	}
}
