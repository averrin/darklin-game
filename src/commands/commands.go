package commands

type Command string

const (
	Say      Command = "/"
	Time     Command = "time"
	Exit     Command = "exit"
	Online   Command = "online"
	Login    Command = "login"
	Help     Command = "help"
	Lookup   Command = "lookup"
	Status   Command = "me"
	Goto     Command = "goto"
	Pick     Command = "pick"
	Drop     Command = "drop"
	Select   Command = "select"
	Unselect Command = "unselect"
	Describe Command = "describe"
	Light    Command = "light"
)

func GetCommands() map[Command][]string {
	return map[Command][]string{
		Say:      []string{},
		Time:     []string{},
		Exit:     []string{},
		Online:   []string{},
		Login:    []string{},
		Help:     []string{},
		Lookup:   []string{},
		Status:   []string{},
		Goto:     []string{},
		Pick:     []string{},
		Drop:     []string{},
		Select:   []string{},
		Unselect: []string{},
		Describe: []string{},
		Light: []string{
			"on", "off",
		},
	}
}
