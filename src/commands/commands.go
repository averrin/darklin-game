package commands

func GetCommands() map[string][]string {
	return map[string][]string{
		"/":      []string{},
		"time":   []string{},
		"exit":   []string{},
		"online": []string{},
		"login":  []string{},
		"help":   []string{},
		"search": []string{},
		"me":     []string{},
		"describe": []string{
			"room",
		},
		"goto": []string{},
		"light": []string{
			"on", "off",
		},
		"pick": []string{},
		"drop": []string{},
	}
}
