package args

func isValidCommand(cmd string) bool {
	validCommands := [...]string{
		"check", "doc", "gen", "help", "ver",
	}

	for i, _ := range validCommands {
		if cmd == validCommands[i] {
			return true
		}
	}

	return false
}

// IsValidDocFlag return true if given flag is valid doc command flag.
func isValidDocFlag(f string) bool {
	validFlags := []string{"-fusesoc", "-no-bold", "-no-config"}

	for i, _ := range validFlags {
		if f == validFlags[i] {
			return true
		}
	}

	return false
}
