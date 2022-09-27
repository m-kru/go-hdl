package utils

import (
	"strings"
)

// hdlGenArgs parses line containing 'hdl:gen' and returns its arguments.
func HdlGenArgs(line []byte) []string {
	splits := strings.Split(string(line), "hdl:gen")
	if splits[1] == "" {
		return nil
	}

	argsSuffix := strings.Trim(splits[1], " \t")
	args := strings.Split(argsSuffix, " ")
	for i := range args {
		args[i] = strings.Trim(args[i], " \t")
	}
	return args
}
