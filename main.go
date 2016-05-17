package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"

	core "./modules/core"
	globalStream "./modules/globalStream"
	player "./modules/player"
	world "./modules/world"

	"gopkg.in/mgo.v2"
)

var VERSION string

func main() {
	log.Println("Core version: " + VERSION)
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

	http.HandleFunc("/ws", gs.GetPlayerHandler(VERSION))

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
