package commands

func GetCommands() map[string][]string {
	return map[string][]string{
		"/":        []string{},
		"time":     []string{},
		"exit":     []string{},
		"online":   []string{},
		"login":    []string{},
		"help":     []string{},
		"lookup":   []string{},
		"me":       []string{},
		"goto":     []string{},
		"pick":     []string{},
		"drop":     []string{},
		"select":   []string{},
		"unselect": []string{},
		"describe": []string{},
		"light": []string{
			"on", "off",
		},
	}
}
