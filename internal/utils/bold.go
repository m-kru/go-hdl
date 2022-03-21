package utils

func BoldCodeTerminal(lang string, code string) string {
	var boldCode string

	switch lang {
	case "vhdl":
		boldCode = VHDLTerminalBold(code)
	default:
		panic("should never happen")
	}

	return boldCode
}
