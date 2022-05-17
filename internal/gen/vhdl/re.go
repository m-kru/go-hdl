package vhdl

import re "regexp"

var thdlGenLine *re.Regexp = re.MustCompile(`\s*--thdl:gen\b`)
