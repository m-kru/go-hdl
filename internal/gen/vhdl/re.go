package vhdl

import re "regexp"

var hdlGenLine *re.Regexp = re.MustCompile(`\s*--hdl:gen\b`)
var hdlFieldArgs *re.Regexp = re.MustCompile(`--hdl: `)

var hdlStartLine *re.Regexp = re.MustCompile(`\s*--hdl:start\b`)
var hdlEndLine *re.Regexp = re.MustCompile(`\s*--hdl:end\b`)
