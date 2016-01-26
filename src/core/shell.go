package core

import (
	"events"
	"log"

	"gopkg.in/readline.v1"
)

// RunShell - interactive mode
func RunShell(stream *chan *events.Event) {
	var completer = readline.NewPrefixCompleter(
		readline.PcItem("time"),
		readline.PcItem("exit"),
		readline.PcItem("online"),
		readline.PcItem("info"),
		readline.PcItem("reset"),
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
		*stream <- events.NewEvent(events.COMMAND, line, "cmd")
	}
}
