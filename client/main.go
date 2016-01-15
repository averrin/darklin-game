package main

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"

	"gopkg.in/readline.v1"
)

func main() {

	var completer = readline.NewPrefixCompleter(
		readline.PcItem("time"),
		readline.PcItem("exit"),
		readline.PcItem("online"),
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

	u := url.URL{Scheme: "ws", Host: "core.darkl.in", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("\n%s", message)
		}
	}()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		println("<< ", line)
		if line == "exit" {
			os.Exit(0)
		}
		err = c.WriteMessage(websocket.TextMessage, []byte(line))
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}
