package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/ugorji/go/codec"

	"gopkg.in/readline.v1"
)

// type EventType int
//
// // Event is atom of event stream
// type Event struct {
// 	Timestamp time.Time
// 	Type      EventType
// 	Payload   interface{}
// 	Sender    string
// }

func runInit(c *websocket.Conn) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	initPath := path.Join(dir, "init.cmd")
	if _, err := os.Stat(initPath); err == nil {
		file, err := os.Open(initPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			log.Println(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				fmt.Println("<<", line)
				c.WriteMessage(websocket.TextMessage, []byte(line))
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

func connect(u url.URL) *websocket.Conn {
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		// log.Fatal("dial:", err)
		time.Sleep(5000 * time.Millisecond)
		return connect(u)
	}
	// log.Println("connected")

	return c
}

func main() {

	var completer = readline.NewPrefixCompleter(
		readline.PcItem("time"),
		readline.PcItem("exit"),
		readline.PcItem("online"),
		readline.PcItem("login"),
		readline.PcItem("help"),
		readline.PcItem("goto",
			readline.PcItem("first"),
			readline.PcItem("second"),
		),
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

	print := func(template string, a ...interface{}) {
		fmt.Fprintf(rl.Stderr(), template+"\n", a...)
	}

	host := flag.String("host", "core.darkl.in", "host of core")
	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *host, Path: "/ws"}
	conn := connect(u)
	defer conn.Close()

	done := make(chan struct{})

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()

	go func() {
		defer conn.Close()
		defer close(done)
		m := 0
		for {
			_, message, err := conn.ReadMessage()
			// log.Println(string(message))
			if err != nil {
				// log.Println("read:", err)
				log.Println(red("Disconnected...") + " wait...")
				time.Sleep(500 * time.Millisecond)
				conn = connect(u)
				continue
				// return
			}
			m++
			var event *Event

			// decoder := json.NewDecoder(bytes.NewReader(message))
			// err = decoder.Decode(&event)
			var mh codec.MsgpackHandle
			dec := codec.NewDecoder(bytes.NewReader(message), &mh)
			err = dec.Decode(&event)
			if err != nil {
				log.Fatal(err)
			}
			switch event.Type {
			case HEARTBEAT:
			case LIGHT:
			case SYSTEMMESSAGE:
				sep := green("|")
				print(sep+" %s", event.Payload)
			case LOGGEDIN:
				sep := green("|")
				print(sep+" %s", event.Payload)
			case ERROR:
				sep := red("!")
				print(sep+" %s", event.Payload)
			case CONNECTED:
				runInit(conn)
			default:
				sep := blue(">")
				// if !strings.HasPrefix(event.Payload.(string), "hi") {
				print(sep+" %s: %s", event.Sender, event.Payload)
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
		// println("<<", line)
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
