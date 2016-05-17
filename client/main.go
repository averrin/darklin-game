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
	"regexp"
	"strings"
	"time"

	commands "../modules/commands"
	events "../modules/events"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/ugorji/go/codec"

	"gopkg.in/readline.v1"
)

//COMMANDS - init commands set
var COMMANDS map[commands.Command][]string

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
	log.Printf("Подключение к %s...", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		// log.Fatal("dial:", err)
		time.Sleep(5000 * time.Millisecond)
		return connect(u)
	}
	// log.Println("connected")

	return c
}

func BuildCompleter(completer *readline.PrefixCompleter) {
	COMMANDS = commands.GetCommands()
	for cmd, children := range COMMANDS {
		item := readline.PcItem(string(cmd))
		for _, child := range children {
			item.Children = append(item.Children, readline.PcItem(child))
		}
		completer.Children = append(completer.Children, item)
	}
}

func ReBuildCompleter(completer *readline.PrefixCompleter, key string, items []interface{}) {
	completer.Children = []*readline.PrefixCompleter{}
	for cmd, children := range COMMANDS {
		item := readline.PcItem(string(cmd))
		if string(cmd) == key {
			COMMANDS[cmd] = []string{}
			for _, child := range items {
				item.Children = append(item.Children, readline.PcItem(string(child.([]byte))))
				COMMANDS[cmd] = append(COMMANDS[cmd], string(child.([]byte)))
			}
		} else {
			for _, child := range children {
				item.Children = append(item.Children, readline.PcItem(child))
			}
		}
		completer.Children = append(completer.Children, item)
	}
}

func main() {

	var completer = readline.NewPrefixCompleter()
	BuildCompleter(completer)
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

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()

	key := regexp.MustCompile(`\[([\w ]+)\]`)
	npc := regexp.MustCompile(`\{([\w ]+)\}`)
	print := func(template string, a ...interface{}) {
		str := fmt.Sprintf(template+"\n", a...)
		str = key.ReplaceAllString(str, "["+green("$1")+"]")
		str = npc.ReplaceAllString(str, "{"+yellow("$1")+"}")
		fmt.Fprint(rl.Stderr(), str)
	}

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
			// log.Println(string(message))
			if err != nil {
				// log.Println("read:", err)
				rl.SetPrompt(">> ")
				rl.Refresh()
				log.Println(red("Отключено...") + " ждем...")
				time.Sleep(500 * time.Millisecond)
				conn = connect(u)
				continue
				// return
			}
			m++
			var event *events.Event

			// decoder := json.NewDecoder(bytes.NewReader(message))
			// err = decoder.Decode(&event)
			var mh codec.MsgpackHandle
			dec := codec.NewDecoder(bytes.NewReader(message), &mh)
			err = dec.Decode(&event)
			if err != nil {
				log.Fatal(err)
			}
			switch event.Type {
			case events.HEARTBEAT:
			case events.LIGHT:
			case events.UNSELECTED:
				rl.SetPrompt(">> ")
				rl.Refresh()
			case events.SELECTED:
				rl.SetPrompt(fmt.Sprintf("> %s > ", green(string(event.Payload.([]byte)))))
				rl.Refresh()
			case events.INTERNALINFO:
				// ii := event.Payload.(actor.InternalInfo)
				if event.Payload != nil {
					ii := event.Payload.(map[interface{}]interface{})
					// log.Println(fmt.Sprintf("%s", ii))
					if string(ii["Type"].([]byte)) == "autocomplete" {
						args := ii["Args"]
						if args != nil {
							ReBuildCompleter(completer, string(ii["Key"].([]byte)), args.([]interface{}))
						}
					}
				}
			case events.ROOMENTER:
				sep := yellow("| ")
				print(sep+"%s %s", event.Sender, event.Payload)
			case events.ROOMEXIT:
				sep := yellow("| ")
				print(sep+"%s %s", event.Sender, event.Payload)
			case events.ROOMCHANGED:
				sep := green("| ")
				print(sep+"%s", event.Payload)
			case events.STATUS:
				sep := blue("| ")
				print(sep+"%s", event.Payload)
			case events.DESCRIBE:
				sep := blue("| ")
				print(sep+"%s", event.Payload)
			case events.SYSTEMMESSAGE:
				sep := green("| ")
				print(sep+"%s", event.Payload)
			case events.LOGGEDIN:
				rl.SetPrompt(">> ")
				rl.Refresh()
				sep := green("| ")
				print(sep+"%s", event.Payload)
			case events.ERROR:
				sep := red("! ")
				print(sep+"%s", event.Payload)
			case events.CONNECTED:
				sep := green("| ")
				print(sep+"Версия сервера: %s", event.Payload)
				runInit(conn)
			default:
				sep := blue("> ")
				// if !strings.HasPrefix(event.Payload.(string), "hi") {
				print(sep+"%s: %s", event.Sender, event.Payload)
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
