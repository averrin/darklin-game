package main

import (
	"actor"
	"events"
	"global"
	"log"
	"player"
	"time"
	"time_stream"

	"gopkg.in/readline.v1"
)

// TestActor just someone who do something
type TestActor struct {
	actor.Actor
}

// ConsumeEvent of couse
func (a TestActor) ConsumeEvent(event events.Event) {
	a.Stream <- event
}

// NewTestActor because i, sucj in golang yet
func NewTestActor(gs chan events.Event) *TestActor {
	a := actor.NewActor("announcer", gs)
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
		case events.SECOND:
			a.SendEvent("global", events.MESSAGE, "Every second, mister")
		case events.MINUTE:
			a.SendEvent("global", events.MESSAGE, "Every minute, boss")
		}
	}
}

func main() {
	gs := global.NewStream()
	stream := gs.Stream
	ts := time_stream.NewStream(stream)
	go ts.Live()

	testActor := NewTestActor(stream)
	go testActor.Live()

	// gs.Subscribe(events.SECOND, testActor)
	gs.Subscribe(events.MINUTE, testActor)

	player := player.NewPlayer(stream)
	// gs.Subscribe(events.MESSAGE, player)
	gs.Streams["player"] = player.Stream
	gs.Streams["time"] = ts.Stream
	go player.Live()

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
			break
		}
		// println("<< ", line)
		stream <- events.Event{time.Now(), events.COMMAND, line, "player"}
	}

}
