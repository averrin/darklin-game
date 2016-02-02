package commands

func GetCommands() map[string][]string {
	return map[string][]string{
		"/":        []string{},
		"time":     []string{},
		"exit":     []string{},
		"online":   []string{},
		"login":    []string{},
		"help":     []string{},
		"search":   []string{},
		"me":       []string{},
		"goto":     []string{},
		"pick":     []string{},
		"drop":     []string{},
		"select":   []string{},
		"unselect": []string{},
		"describe": []string{
			"room",
		},
		"light": []string{
			"on", "off",
		},
	}
}
