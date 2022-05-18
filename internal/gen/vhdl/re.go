package vhdl

import re "regexp"

var thdlGenLine *re.Regexp = re.MustCompile(`\s*--thdl:gen\b`)

var thdlStartLine *re.Regexp = re.MustCompile(`\s*--thdl:start\b`)
var thdlEndLine *re.Regexp = re.MustCompile(`\s*--thdl:end\b`)
