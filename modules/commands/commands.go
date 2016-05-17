package commands

// Command type
type Command string

// Commands list
var Commands []Command

// NewCommand - constructor
func NewCommand(s string) Command {
	c := Command(s)
	Commands = append(Commands, c)
	return c
}

var (
	Say      = NewCommand("/")
	Time     = NewCommand("time")
	Exit     = NewCommand("exit")
	Online   = NewCommand("online")
	Login    = NewCommand("login")
	Help     = NewCommand("help")
	Lookup   = NewCommand("lookup")
	Status   = NewCommand("me")
	Goto     = NewCommand("goto")
	Pick     = NewCommand("pick")
	Drop     = NewCommand("drop")
	Select   = NewCommand("select")
	Unselect = NewCommand("unselect")
	Describe = NewCommand("describe")
	Light    = NewCommand("light")
	Use      = NewCommand("use")
)

func GetCommands() map[Command][]string {
	c := map[Command][]string{}
	for _, cmd := range Commands {
		c[cmd] = []string{}
	}
	c[Light] = []string{
		"on", "off",
	}
	return c
}
