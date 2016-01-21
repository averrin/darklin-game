package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"gopkg.in/readline.v1"
)

func connect(u url.URL) *websocket.Conn {
	// log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		// log.Fatal("dial:", err)
		time.Sleep(5000 * time.Millisecond)
		return connect(u)
	}
	// log.Println("connected")
	return c
}

//Event - atomic action
type Event struct {
	Timestamp time.Time
	Type      int
	Payload   interface{}
	Sender    string
}

func main() {
	host := flag.String("host", "core.darkl.in", "host of core")
	count := flag.Int("count", 500, "connections count")
	delay := flag.Int("delay", 500, "message delay")
	flag.Parse()

	var completer = readline.NewPrefixCompleter(
		readline.PcItem("time"),
		readline.PcItem("exit"),
		readline.PcItem("online"),
		readline.PcItem("login"),
	)
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       fmt.Sprintf(">stress [c%vd%vms]> ", *count, *delay),
		HistoryFile:  "/tmp/readline.tmp",
		AutoComplete: completer,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	log.SetOutput(rl.Stderr())
	log.SetPrefix("")

	u := url.URL{Scheme: "ws", Host: *host, Path: "/ws"}
	// var connections []*websocket.Conn
	// var conn *websocket.Conn
	lg := 0
	sended := 0
	dc := 0
	// failed := 0
	for index := 0; index < *count; index++ {
		go func(index int) {
			conn := connect(u)
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("login %v 123", index)))
			if index%2 == 0 {
				conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("goto second", index)))
			}
			go func(index int, conn *websocket.Conn) {
				defer conn.Close()
				for {
					conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("/hi from %v", index)))
					sended++
					time.Sleep(time.Duration(*delay * int(time.Millisecond)))
				}
			}(index, conn)
			go func(conn *websocket.Conn) {
				defer conn.Close()
				var event *Event
				for {
					_, message, err := conn.ReadMessage()
					if err != nil {
						log.Println("Disconnected")
						log.Printf("\nusers: %v/%v/%v", lg, dc, *count)
						dc++
						break
					}
					decoder := json.NewDecoder(bytes.NewReader(message))
					decoder.Decode(&event)
					switch event.Type {
					case 14:
						lg++
						log.Printf("\nusers: %v/%v/%v", lg, dc, *count)
						// case 12:
						// 	failed++
						// 	log.Printf("\nmessages: %v/%v", failed, sended)
					}
				}
			}(conn)
			// connections = append(connections, conn)
		}(index)
		time.Sleep(time.Duration(50 * int(time.Millisecond)))
	}

	// done := make(chan struct{})

	// go func() {
	// 	defer conn.Close()
	// 	defer close(done)
	// 	for {
	// 		_, message, err := conn.ReadMessage()
	// 		if err != nil {
	// 			log.Println("read:", err)
	// 			log.Println("Disconnected... wait...")
	// 			time.Sleep(500 * time.Millisecond)
	// 			conn = connect(u)
	// 			continue
	// 			// return
	// 		}
	// 		var event *Event
	// 		decoder := json.NewDecoder(bytes.NewReader(message))
	// 		err = decoder.Decode(&event)
	// 		switch event.Type {
	// 		case 8:
	// 		default:
	// 			log.Printf("\n%s: %v", event.Sender, event.Payload)
	// 		}
	// 	}
	// }()

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
		// err = conn.WriteMessage(websocket.TextMessage, []byte(line))
		// if err != nil {
		// 	log.Println("write:", err)
		// 	return
		// }
	}
}
