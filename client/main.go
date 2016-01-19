package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"gopkg.in/readline.v1"
)

func connect(u url.URL) *websocket.Conn {
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		// log.Fatal("dial:", err)
		time.Sleep(5000 * time.Millisecond)
		return connect(u)
	}
	log.Println("connected")
	return c
}

//Event - event from server
type Event struct {
	Timestamp time.Time
	Type      int
	Payload   interface{}
	Sender    string
}

func main() {

	var completer = readline.NewPrefixCompleter(
		readline.PcItem("time"),
		readline.PcItem("exit"),
		readline.PcItem("online"),
		readline.PcItem("login"),
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

	host := flag.String("host", "core.darkl.in", "host of core")
	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *host, Path: "/ws"}
	conn := connect(u)
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer conn.Close()
		defer close(done)
		m := 0
		for {
			_, message, err := conn.ReadMessage()
			log.Println(string(message))
			if err != nil {
				// log.Println("read:", err)
				log.Println("Disconnected... wait...")
				time.Sleep(500 * time.Millisecond)
				conn = connect(u)
				continue
				// return
			}
			m++
			var event *Event
			decoder := json.NewDecoder(bytes.NewReader(message))
			err = decoder.Decode(&event)
			switch event.Type {
			case 8:
			default:
				// if !strings.HasPrefix(event.Payload.(string), "hi") {
				// log.Printf("\n%s: %v", event.Sender, event.Payload)
				// }
			}
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
		err = conn.WriteMessage(websocket.TextMessage, []byte(line))
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}
