package main

import (
	"log"
	"net/http"
	"time"

	"gopkg.in/readline.v1"
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
	gs := NewGlobalStream()
	stream := gs.Stream
	ts := NewTimeStream(stream)
	go ts.Live()

	testActor := NewTestActor(stream)
	go testActor.Live()

	gs.Subscribe(SECOND, testActor)
	gs.Subscribe(MINUTE, testActor)

	gs.Streams["time"] = ts.Stream

	http.HandleFunc("/ws", gs.CmdHandler)
	log.Println("Serving")
	// http.Handle("/", http.FileServer(http.Dir(".")))
	go http.ListenAndServe(":80", nil)
	go gs.Live()

	var completer = readline.NewPrefixCompleter(
		readline.PcItem("time"),
		readline.PcItem("exit"),
	)
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       ">> ",
		HistoryFile:  "/tmp/readline.tmp",
		AutoComplete: completer,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	log.SetOutput(rl.Stderr())
	log.SetPrefix("")

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			log.Fatal(err)
			break
		}
		if line == "" {
			continue
		}
		println("<< ", line)
		stream <- Event{time.Now(), COMMAND, line, "cmd"}
	}
}
