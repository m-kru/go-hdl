package vhdl

import re "regexp"

var thdlGenLine *re.Regexp = re.MustCompile(`\s*--thdl:gen\b`)
var thdlFieldArgs *re.Regexp = re.MustCompile(`--thdl: `)

var thdlStartLine *re.Regexp = re.MustCompile(`\s*--thdl:start\b`)
var thdlEndLine *re.Regexp = re.MustCompile(`\s*--thdl:end\b`)
