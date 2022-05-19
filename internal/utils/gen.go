package utils

import (
	"strings"
)

// thdlGenArgs parses line containing 'thdl:gen' and returns its arguments.
func ThdlGenArgs(line []byte) []string {
	splits := strings.Split(string(line), "thdl:gen")
	if splits[1] == "" {
		return nil
	}

	argsSuffix := strings.Trim(splits[1], " \t")
	args := strings.Split(argsSuffix, " ")
	for i, _ := range args {
		args[i] = strings.Trim(args[i], " \t")
	}
	return args
}
