package args

func isValidCommand(cmd string) bool {
	validCommands := [...]string{
		"doc", "gen", "help", "ver", "vet",
	}

	for i := range validCommands {
		if cmd == validCommands[i] {
			return true
		}
	}

	return false
}
